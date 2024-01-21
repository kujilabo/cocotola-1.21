import { useEffect, useRef, useState } from 'react';

import { Container, Breadcrumb, BreadcrumbItem, BreadcrumbLink } from '@chakra-ui/react';
import { Howl, SoundSpriteDefinitions } from 'howler';
import { useParams, Link } from 'react-router-dom';

import { useWorkbookRetrieveStore } from '../api/workbook_retrieve';

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
export const WorkbookView = (): JSX.Element => {
  const { _workbookId } = useParams();
  const once = useRef(false);
  const workbooks = useWorkbookRetrieveStore((state) => state.workbooks);
  const state = useWorkbookRetrieveStore((state) => state.state);
  const retrieveWorkbook = useWorkbookRetrieveStore((state) => state.retrieveWorkbook);
  const [status, setStatus] = useState('IDLE');
  const [soundId, setSoundId] = useState(0);
  const [trackIndex, setTrackIndex] = useState(0);
  const [loop, setLoop] = useState(false);
  const [trackLength, setTrackLength] = useState(0);
  const tracks = [1, 2, 3, 4];
  const [player, setPlayer] = useState(0);

  const [sprite, setSprite] = useState<SoundSpriteDefinitions>({});
  const [src, setSrc] = useState('');

  const changePlayer = () => {
    if (player === 0) {
      setPlayer(1);
    } else {
      setPlayer(0);
    }
  };

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
      // onend: (id) => {
      //   console.log('サウンド終了!!', id);
      //   setStatus('IDLE');
      // },
    });
  };

  // src and wait(dst)
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
    for (let i = 0; i < workbook.problems.length; i++) {
      const problem = workbook.problems[i];
      src += problem.properties['srcAudioContent'];
      totalAudioLength += +problem.properties['srcAudioLength'];
    }

    let offset = 0;
    const tmpSprite: SoundSpriteDefinitions = {};
    if (player === 0) {
      for (let i = 0; i < workbook.problems.length; i++) {
        const problem = workbook.problems[i];
        const audioLength = +problem.properties['srcAudioLength'];
        const trackNo = makeTrackNo(i + 1);
        tmpSprite[trackNo] = [offset, audioLength];
        offset += audioLength;
      }
      setTrackLength(workbook.problems.length);
    } else if (player === 1) {
      for (let i = 0; i < workbook.problems.length; i++) {
        const problem = workbook.problems[i];
        const audioLength = +problem.properties['srcAudioLength'];
        const trackNo1 = makeTrackNo(i * 2 + 1);
        tmpSprite[trackNo1] = [offset, audioLength];
        offset += audioLength;

        const trackNo2 = makeTrackNo(i * 2 + 2);
        tmpSprite[trackNo2] = [totalAudioLength, 500];
      }
      setTrackLength(workbook.problems.length * 2);
    }

    setSrc(src);
    setSprite(tmpSprite);
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

  const palyButton = () => {
    switch (status) {
      case 'PLAYING':
        return <button onClick={() => pause()}>Pause</button>;
      case 'IDLE':
        return <button onClick={() => start()}>Play</button>;
      case 'PAUSE':
        return <button onClick={() => resume()}>Resume</button>;
    }
  };
  const stopButton = () => {
    switch (status) {
      case 'PLAYING':
        return <button onClick={() => stop()}>stop</button>;
      default:
        return <> </>;
    }
  };

  let main = <></>;
  if (workbookId !== 0 && workbookId in workbooks) {
    main = (
      <div>
        <button onClick={() => changePlayer()}>change player</button>
        {player === 0 ? <div>Source only</div> : <div>Source and Destination</div>}
        <br />
        {palyButton()}
        <br />
        <div>
          {tracks.map((track) => {
            return (
              <div key={track}>
                <button onClick={() => sound.play('track0' + ((track - 1) * 2 + 1).toString())}>
                  {track}
                </button>
              </div>
            );
          })}
        </div>
        <br />
        <button onClick={() => setLoop(!loop)}>{loop ? <>LOOP</> : <>NOT LOOP</>}</button>
        <br />
        {stopButton()}
        <br />
      </div>
    );
  }

  return (
    <Container>
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
    </Container>
  );
};
