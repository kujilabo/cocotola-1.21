import { Navigate } from 'react-router';
import { useEffect, useState } from 'react';

import logo from '@/assets/react.svg';
import { Head } from '@/components/head';
import { LoginForm } from '@/features/auth/components/LoginForm';
import { useAuthStore } from '@/stores/auth';
import { useGlobalAudioPlayer } from 'react-use-audio-player';

import { Howl, Howler } from 'howler';
import { header, base64_weather, base64_fortune } from './data';
export const Landing = () => {
  const getUserInfo = useAuthStore((state) => state.getUserInfo);
  const [play, setPlay] = useState(false);
  const [soundId, setSoundId] = useState(0);
  const userInfo = getUserInfo();
  let tmpSoundId = 0;
  let tmpSoundIndex = 0;
  console.log('AAAAAAAAAAAAAAAAAAAAAAAAAaa');
  // let index = 0;
  // let soundId = 0;
  const { load } = useGlobalAudioPlayer();
  const [index, setIndex] = useState(0);
  // console.log('index', index);
  // if (userInfo) {
  //   return <Navigate to="/app" />;
  // }

  let snd_3 = new Howl({
    src: [header + ',' + base64_fortune + base64_weather],
    sprite: {
      track01: [0, 2736],
      track02: [2736 + 3648, 2000],
      track03: [2736, 3648],
    },
    loop: false, // 繰り返し
    volume: 1.0, // 音量
    rate: 1.0, // 再生速度
    onplay: (id) => {
      setPlay(true);
      console.log('サウンド再生!!', id);
      tmpSoundId = id;
      setSoundId(id);
    },
    onstop: () => {
      console.log('サウンド停止!!');
    },
    onpause: (id) => {
      console.log('サウンド一時停止!!', id);
    },
  });
  const resume = () => {
    console.log('soundId==>', soundId);
    sound.play(soundId);
  };

  useEffect(() => {
    console.log('index', index);
    if (index !== 0) {
      const trackNo = 'track0' + index.toString();
      console.log('trackNo', trackNo);
      tmpSoundIndex = index;
      sound.once('end', function (id: number) {
        console.log('サウンド終了!!', id);
        console.log('index', index);
        if (index >= 3) {
          setPlay(false);
          return;
        }
        setIndex(index + 1);
      });
      sound.play(trackNo);
    }
  }, [index]);
  let snd_1 = new Howl({
    src: header + ',' + base64_fortune,
    loop: false, // 繰り返し
    volume: 1.0, // 音量
    rate: 1.0, // 再生速度
    onplay: () => {
      console.log('サウンド再生!!');
    },
    onstop: () => {
      console.log('サウンド停止!!');
    },
    onpause: () => {
      console.log('サウンド一時停止!!');
    },
    onend: () => {
      console.log('サウンド終了!!');
    },
  });
  let snd_2 = new Howl({
    src: header + ',' + base64_weather,
    loop: false, // 繰り返し
    volume: 1.0, // 音量
    rate: 1.0, // 再生速度
    onplay: () => {
      console.log('サウンド再生!!');
    },
    onstop: () => {
      console.log('サウンド停止!!');
    },
    onpause: () => {
      console.log('サウンド一時停止!!');
    },
    onend: () => {
      console.log('サウンド終了!!');
    },
  });
  // load(base64_1);
  const [sound, setSound] = useState(snd_3);

  const pause = () => {
    console.log('soundId==>', soundId);
    sound.pause(soundId);
  };

  return (
    <>
      <Head description="Welcome to bulletproof react" />
      <div className="bg-white h-[100vh] flex items-center">
        <div className="max-w-7xl mx-auto text-center py-12 px-4 sm:px-6 lg:py-16 lg:px-8">
          <h2 className="text-3xl font-extrabold tracking-tight text-gray-900 sm:text-4xl">
            <span className="block">Bulletproof React</span>
          </h2>
          <img src={logo} alt="react" />
          <LoginForm />
          Landing
          <button onClick={() => snd_1.play()}>PLAY</button>
          <button onClick={() => snd_2.play()}>PLAY</button>
          <button
            onClick={() => {
              // snd_3.play('track01');
              // snd_3.play('track02');
              // index = 1;
              tmpSoundId = snd_3.play('track01');
              // snd_3.pause();
              console.log('snd3 play', snd_3);
              console.log('soundId', tmpSoundId);
              // snd_3.play(0);
            }}
          >
            PLAY
          </button>
          <br />
          <button onClick={() => snd_1.play()}>PLAY</button>
          <button onClick={() => snd_2.play()}>PLAY</button>
          <button
            onClick={() => {
              setIndex(1);
            }}
          >
            PLAY
          </button>
          <button
            onClick={() => {
              setIndex(1);
            }}
          >
            PLAY
          </button>
          <br />
          <button onClick={pause}>pause</button>
          <button onClick={() => resume()}>saikai</button>
          <br />
          {play ? 'Playing' : 'Idle'}
        </div>
      </div>
    </>
  );
};
