import { Navigate } from 'react-router';
import { useEffect, useState } from 'react';

import logo from '@/assets/react.svg';
import { Head } from '@/components/head';
import { LoginForm } from '@/features/auth/components/LoginForm';
import { useAuthStore } from '@/stores/auth';

import { Howl, Howler } from 'howler';
import { header, base64_weather, base64_fortune } from './data';
export const Landing = () => {
  const getUserInfo = useAuthStore((state) => state.getUserInfo);
  const [status, setStatus] = useState('IDLE');
  const [soundId, setSoundId] = useState(0);
  const [loop, setLoop] = useState(false);
  const userInfo = getUserInfo();
  console.log('AAAAAAAAAAAAAAAAAAAAAAAAAaa', status);
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

  const start = () => {
    setIndex(1);
  };
  const stop = () => {
    console.log('soundId', soundId);
    sound.stop();
    setStatus('IDLE');
    setIndex(0);
  };
  const resume = () => {
    sound.play(soundId);
  };
  const pause = () => {
    sound.pause(soundId);
  };

  useEffect(() => {
    if (index !== 0) {
      const trackNo = 'track0' + index.toString();
      console.log('trackNo', trackNo);
      console.log('soundId', soundId);
      sound.once('end', function (id: number) {
        console.log('サウンド終了!!', id);
        console.log('index', index);
        if (index >= 3) {
          if (loop) {
            start();
          } else {
            stop();
          }
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

  const [sound, setSound] = useState(snd_3);

  const palyButton = () => {
    switch (status) {
      case 'PLAYING':
        return <button onClick={() => pause()}>pause</button>;
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
          <br />
          <button onClick={() => setLoop(!loop)}>{loop ? <>LOOP</> : <>NOT LOOP</>}</button>
          <br />
          <br />
          {palyButton()}
          <br />
          {stopButton()}
          <br />
          <br />
          <h1>{status}</h1>
          <h1>{index}</h1>
        </div>
      </div>
    </>
  );
};
