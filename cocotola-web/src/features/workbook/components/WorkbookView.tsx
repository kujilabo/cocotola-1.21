import { useEffect, useLayoutEffect, useRef, useState } from 'react';

import { ChevronDownIcon } from '@chakra-ui/icons';
import {
  Box,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  Button,
  ButtonGroup,
  Card,
  CardBody,
  Container,
  Flex,
  Heading,
  HStack,
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Spacer,
  Text,
} from '@chakra-ui/react';
import { Howl, SoundSpriteDefinitions } from 'howler';
import { FaPlay } from 'react-icons/fa';
import { RxLoop } from 'react-icons/rx';
import { useParams, Link } from 'react-router-dom';

import { BasicButton } from '@/components/buttons/BasicButton';

import { useWorkbookRetrieveStore } from '../api/workbook_retrieve';
import { EnglishSentence } from '../types';

export const header = 'data:audio/mp3;base64';
const toNumber = (s: string | undefined): number => {
  if (s === undefined) {
    return 0;
  }
  return parseInt(s);
};
const pad = (num: number, size: number): string => {
  let s = num + '';
  while (s.length < size) s = '0' + s;
  return s;
};
const makeTrackNo = (num: number): string => {
  return 'track' + pad(num, 2);
};
const createSpriteForSourceOnly = (englishSentences: EnglishSentence[]): SoundSpriteDefinitions => {
  let offset = 0;
  const tmpSprite: SoundSpriteDefinitions = {};
  for (let i = 0; i < englishSentences.length; i++) {
    const englishSentence = englishSentences[i];
    const audioLength = englishSentence.srcAudioLengthMillisecond;
    const trackNo = makeTrackNo(i + 1);
    tmpSprite[trackNo] = [offset, audioLength];
    offset += audioLength;
  }
  console.log('tmpSprite', tmpSprite);
  return tmpSprite;
};
const createSpriteForSourceAndDestination = (
  englishSentences: EnglishSentence[],
  totalAudioLength: number
): SoundSpriteDefinitions => {
  let offset = 0;
  const tmpSprite: SoundSpriteDefinitions = {};

  for (let i = 0; i < englishSentences.length; i++) {
    const englishSentence = englishSentences[i];
    const audioLength = englishSentence.srcAudioLengthMillisecond;
    const trackNo1 = makeTrackNo(i * 2 + 1);
    tmpSprite[trackNo1] = [offset, audioLength];
    offset += audioLength;

    const trackNo2 = makeTrackNo(i * 2 + 2);
    tmpSprite[trackNo2] = [totalAudioLength, 500];
  }
  console.log('tmpSprite', tmpSprite);
  return tmpSprite;
};

const createSprite = (
  player: number,
  englishSentences: EnglishSentence[],
  totalAudioLength: number
): SoundSpriteDefinitions => {
  if (player === PLAYER_SOURCE_ONLY) {
    return createSpriteForSourceOnly(englishSentences);
  } else if (player === PLAYER_SOURCE_AND_DESTINATION) {
    return createSpriteForSourceAndDestination(englishSentences, totalAudioLength);
  } else {
    return {};
  }
};
const createCard = (
  englishSentences: EnglishSentence[],
  activeIndex: number,
  scrollRef: React.RefObject<HTMLDivElement>
) => {
  const activeColor = 'green.300';
  const inactiveColor = 'gray.100';
  return (
    <Box overflow="scroll" height="calc(100vh - 250px)">
      {englishSentences.map((englishSentence: EnglishSentence, i: number) => {
        const color = activeIndex == i ? activeColor : inactiveColor;
        const ref = activeIndex == i || (i == 0 && activeIndex < 0) ? scrollRef : null;
        return (
          <Box key={i} ref={ref}>
            <Box p={1}>
              <HStack>
                <Spacer />
                <Card bg="gray.100" width="90%" borderColor={color} borderWidth={1}>
                  <CardBody>
                    <Box>
                      <Heading size="xs">{englishSentence.dstText}</Heading>
                      <Text pt="2" fontSize="sm">
                        {englishSentence.srcText}
                      </Text>
                    </Box>
                  </CardBody>
                </Card>
              </HStack>
            </Box>
            <Box p={1}>
              <HStack>
                <Spacer />
                <Box>
                  <ButtonGroup>
                    <Button>Learned!</Button>
                    <Button>PLay!</Button>
                  </ButtonGroup>
                </Box>
              </HStack>
            </Box>
          </Box>
        );
      })}
    </Box>
  );
};
const PLAYER_SOURCE_ONLY = 0;
const PLAYER_SOURCE_AND_DESTINATION = 1;
export const WorkbookView = (): JSX.Element => {
  const { _workbookId } = useParams();
  const once = useRef(false);
  const scrollRef = useRef<HTMLDivElement>(null);
  const workbooks = useWorkbookRetrieveStore((state) => state.workbooks);
  const state = useWorkbookRetrieveStore((state) => state.state);
  const retrieveWorkbook = useWorkbookRetrieveStore((state) => state.retrieveWorkbook);
  const [status, setStatus] = useState('IDLE');
  const [soundId, setSoundId] = useState(0);
  const [trackIndex, setTrackIndex] = useState(0);
  const [loop, setLoop] = useState(false);
  const [trackLength, setTrackLength] = useState(0);
  const [player, setPlayer] = useState(0);

  const [sprite, setSprite] = useState<SoundSpriteDefinitions>({});
  const [src, setSrc] = useState('');

  const newSound = (src: string, sprite: SoundSpriteDefinitions): Howl => {
    return new Howl({
      src: [src],
      sprite: sprite,
      loop: false, // 繰り返し
      volume: 1.0, // 音量
      rate: 1.0, // 再生速度
      onplay: (id) => {
        console.log('サウンド再生!!', id);
        setStatus('PLAYING');
        setSoundId(id);
      },
      onstop: () => {
        console.log('サウンド停止!!');
      },
      onpause: (id) => {
        console.log('サウンド一時停止!!', id);
        setStatus('PAUSE');
      },
    });
  };

  const [sound, setSound] = useState(newSound(src, sprite));

  const start = () => {
    setTrackIndex(1);
  };
  const stop = () => {
    sound.stop();
    setStatus('IDLE');
    setTrackIndex(0);
  };
  const resume = () => {
    sound.play(soundId);
  };
  const pause = () => {
    sound.pause(soundId);
  };

  const workbookId = toNumber(_workbookId);
  useEffect(() => {
    console.log('init src and sprite');
    if (workbookId === 0) {
      return;
    }
    if (!(workbookId in workbooks)) {
      return;
    }
    const workbook = workbooks[workbookId];
    let src = header + ',';
    let totalAudioLength = 0;
    for (let i = 0; i < workbook.englishSentences.sentences.length; i++) {
      const sentence = workbook.englishSentences.sentences[i];
      src += sentence.srcAudioContent;
      totalAudioLength += +sentence.srcAudioLengthMillisecond;
    }

    const tmpSprite = createSprite(player, workbook.englishSentences.sentences, totalAudioLength);
    setSprite(tmpSprite);
    setTrackLength(Object.keys(tmpSprite).length);
    setSrc(src);
    setSound(newSound(src, tmpSprite));
  }, [workbookId, workbooks, player]);

  useEffect(() => {
    if (once.current === false) {
      once.current = true;
      if (!(workbookId in workbooks) && state === 'idle') {
        const f = async () => {
          await retrieveWorkbook(1);
        };
        f().catch(console.error);
      }
    }
  }, [workbookId, workbooks, state, retrieveWorkbook]);

  useLayoutEffect(() => {
    scrollRef.current?.scrollIntoView();
  }, [scrollRef, trackIndex]);

  // play track
  useEffect(() => {
    if (trackIndex !== 0) {
      const trackNo = makeTrackNo(trackIndex);
      sound.once('end', function (id: number) {
        console.log('サウンド終了!!', id, trackNo);
        if (trackIndex >= trackLength) {
          if (loop) {
            start();
          } else {
            stop();
          }
          return;
        }
        setTrackIndex(trackIndex + 1);
      });
      sound.play(trackNo);
    }
  }, [trackIndex]);

  const playButton = () => {
    switch (status) {
      case 'PLAYING':
        return <BasicButton onClick={() => pause()} value="PAUSE" />;
      case 'IDLE':
        // return <BasicButton onClick={() => start()} value="PLAY" />;
        return (
          <IconButton
            colorScheme="teal"
            aria-label="Play"
            size="lg"
            icon={<FaPlay />}
            onClick={() => start()}
          />
        );
      case 'PAUSE':
        return <BasicButton onClick={() => resume()} value="RESUME" />;
    }
  };
  const stopButton = () => {
    // switch (status) {
    //   case 'PLAYING':
    //     return <BasicButton onClick={() => stop()} value="STOP" />;
    //   default:
    //     return <> </>;
    // }
    return <></>;
  };
  if (workbookId === 0) {
    return <> </>;
  }
  if (!(workbookId in workbooks)) {
    return <> </>;
  }
  const workbook = workbooks[workbookId];

  console.log('trackIndex - 1', trackIndex - 1);
  const card = createCard(workbook.englishSentences.sentences, trackIndex - 1, scrollRef);

  const footer = (
    // <ButtonGroup>
    <Flex>
      <Spacer />
      <ButtonGroup>
        <IconButton
          colorScheme="teal"
          variant={loop ? 'solid' : 'outline'}
          aria-label="Loop"
          size="lg"
          icon={<RxLoop />}
          onClick={() => setLoop(!loop)}
        />
        <Menu>
          <MenuButton colorScheme="teal" size="lg" as={Button} rightIcon={<ChevronDownIcon />}>
            {player === PLAYER_SOURCE_ONLY ? <>Source Only</> : <>Source and Destination</>}
          </MenuButton>
          <MenuList>
            <MenuItem onClick={() => setPlayer(PLAYER_SOURCE_ONLY)}>Source Only</MenuItem>
            <MenuItem onClick={() => setPlayer(PLAYER_SOURCE_AND_DESTINATION)}>
              Source and Destination
            </MenuItem>
          </MenuList>
        </Menu>
        {playButton()}{' '}
      </ButtonGroup>
    </Flex>
    //
  );

  //   <div>
  //   {tracks.map((track) => {
  //     return (
  //       <div key={track}>
  //         <button onClick={() => sound.play('track0' + ((track - 1) * 2 + 1).toString())}>
  //           {track}
  //         </button>
  //       </div>
  //     );
  //   })}
  // </div>
  let main = <></>;
  if (workbookId !== 0 && workbookId in workbooks) {
    main = (
      <Box flex="1">
        <br />
        <button onClick={() => setLoop(!loop)}>{loop ? <>LOOP</> : <>NOT LOOP</>}</button>
        <br />
        {stopButton()}
        <br />
        {card}
        <br />
      </Box>
    );
  }

  return (
    <Container bg="white" minH="calc(100vh - 72px)">
      {/* <Box> */}
      <Box height="calc(100vh - 72px - 60px)" bg="orange.100">
        <Breadcrumb>
          <BreadcrumbItem>
            <BreadcrumbLink as={Link} to="/">
              Home
            </BreadcrumbLink>
          </BreadcrumbItem>

          <BreadcrumbItem>
            <BreadcrumbLink href="#">Docs</BreadcrumbLink>
          </BreadcrumbItem>

          <BreadcrumbItem isCurrentPage>
            <BreadcrumbLink href="#">Breadcrumb</BreadcrumbLink>
          </BreadcrumbItem>
        </Breadcrumb>
        {main}
      </Box>
      <Box p="2" bg="tomato" height="60px">
        {footer}
      </Box>
    </Container>
  );
};
