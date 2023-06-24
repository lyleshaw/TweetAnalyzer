import { Container, Card, Link, Row, Col, Text, Spacer, Input, Modal, Button, styled, Loading, Dropdown, Grid } from "@nextui-org/react";
import React, { useId, Suspense, useCallback, useState, useTransition } from 'react'
import { atomWithObservable } from 'jotai/utils'
import { useAtom, useAtomValue, useSetAtom } from 'jotai/react'
import { atom } from 'jotai/vanilla'
import { Observable } from 'rxjs'
import { FormElement } from '@nextui-org/react/types/input/input-props'
import { Model, ModelChange } from "./components/model-change";
import { GrPowerCycle } from 'react-icons/gr'

export const SendButton = styled('button', {
  // reset button styles
  background: 'transparent',
  border: 'none',
  padding: 0,
  // styles
  width: '24px',
  margin: '0 10px',
  dflex: 'center',
  bg: '$primary',
  borderRadius: '$rounded',
  cursor: 'pointer',
  transition: 'opacity 0.25s ease 0s, transform 0.25s ease 0s',
  svg: {
    size: '100%',
    padding: '4px',
    transition: 'transform 0.25s ease 0s, opacity 200ms ease-in-out 50ms',
    boxShadow: '0 5px 20px -5px rgba(0, 0, 0, 0.1)',
  },
  '&:hover': {
    opacity: 0.8
  },
  '&:active': {
    transform: 'scale(0.9)',
    svg: {
      transform: translate(24px, -24px);
      opacity: 0;
    }
  }
`;

type SendIconProps = {
  fill?: string;
  filled?: boolean;
  size?: number;
  height?: number;
  width?: number;
  label?: string;
  className?: string;
};

const SendIcon = ({
  fill = "currentColor",
  filled,
  size,
  height,
  width,
  label,
  className,
}: SendIconProps) => {
  return (
    <svg
      data-name="Iconly/Curved/Lock"
      xmlns="http://www.w3.org/2000/svg"
      width={size || width || 24}
      height={size || height || 24}
      viewBox="0 0 24 24"
      className={className}
    >
      <g transform="translate(2 2)">
        <path
          d="M19.435.582A1.933,1.933,0,0,0,17.5.079L1.408,4.76A1.919,1.919,0,0,0,.024,6.281a2.253,2.253,0,0,0,1,2.1L6.06,11.477a1.3,1.3,0,0,0,1.61-.193l5.763-5.8a.734.734,0,0,1,1.06,0,.763.763,0,0,1,0,1.067l-5.773,5.8a1.324,1.324,0,0,0-.193,1.619L11.6,19.054A1.91,1.91,0,0,0,13.263,20a2.078,2.078,0,0,0,.25-.01A1.95,1.95,0,0,0,15.144,18.6L19.916,2.525a1.964,1.964,0,0,0-.48-1.943"
          fill={fill}
        />
      </g>
    </svg>
  );
};

const twitterIdAtom = atom<string | null>(null);

const contentQueryRequestIdStateAtom = atom<number>(0)

const models:Record<string, string> = {
  'Open-AI': 'https://weet-api-boe.aireview.tech/api/get_tweet_analysis',
  'Claude': 'https://tweet-api.aireview.tech/api/get_tweet_analysis'
}

const modelAtom = atom('Open-AI');


const contentAtom = atomWithObservable((get) => {
  const id = get(twitterIdAtom)
  const mode = get(modelAtom);
  const api = models[mode];
  return new Observable<string | null>((subscriber) => {
    const abortController = new AbortController();
    if (id === null) {
      // no value
      subscriber.next(null);
    } else {
      async function fetchData() {
        const response = await fetch(`${api}?twitter_id=${id}`, {
          headers: {
            'Content-Type': 'application/json',
          },
          signal: abortController.signal,
          method: 'GET',
        })

        if (!response.ok) {
          const error = await response.json();
          console.error(error.error);
          throw new Error("Request failed");
        }
        const data = response.body;
        if (!data) {
          throw new Error("No data");
        }

        const reader = data.getReader();
        const decoder = new TextDecoder("utf-8");
        let done = false;

        while (!done) {
          const { value, done: readerDone } = await reader.read();
          if (value) {
            const char = decoder.decode(value);

            if (char.startsWith("Data:")) {
              subscriber.next(char.substring(5));
            } else {
              subscriber.next(char);
            }
          }
          done = readerDone;
        }
      }

      fetchData().catch((err) => subscriber.error(err))
    }

    return () => {
      abortController.abort();
    };
  });
});

const Content = () => {
  const content = useAtomValue(contentAtom);
  return <>{content}</>;
};

export default function App() {
  const [visible, setVisible] = useAtom(visibleAtom)
  const [model, setModel] = useAtom(modelAtom)
  const handler = useCallback(() => setVisible(true), [setVisible]);
  const closeHandler = useCallback(() => {
    setVisible(false);
    console.log("closed");
  }, [setVisible]);
  const [isLoading, startTransition] = useTransition()
  const [twitterId, setTwitterId] = useAtom(twitterIdAtom)
  const setRequestId = useSetAtom(contentQueryRequestIdStateAtom)
  const handler = useCallback(() => {
    modalRef.current?.showModal();
  }, [modalRef]);

  const handleConfirm = async () => {
    startTransition(() => {
      if (input !== "") {
        setTwitterId(input);
      }
    });
  };

  const onSelect = (mode: Model)=>{
    setModel(Object.keys(mode)[0])
  }

  return (
    <Container sm display="flex" gap={7}
      css={{
        marginTop: '4em'
      }}
    >
      <Col css={{
        textAlign: 'center'
      }}>
        <Text
          transform="full-size-kana"
          h1
          size={50}
          css={ {
            textGradient: "to right, #006E3A 8%, #166BB5 100%",
            '@xsMax': {fontSize: '2.5rem'},
          } }
          weight="bold"
        >
          推文分析器
        </Text>
        <Text
          transform="full-size-kana"
          h1
          size={ 37 }
          css={ {
            textGradient: "to right, #006E3A 8%, #166BB5 100%",
          } }
          weight="bold"
        >
          Tweet Analyzer
        </Text>
        <ModelChange models={models} defaultSelect={model} onSelect={onSelect}>
          <Text
            transform="full-size-kana"
            h1
            size={24}
            css={ {
              textGradient: "to right, #006E3A 8%, #166BB5 100%",
              '@xsMax': {fontSize: '2.5rem'},
            } }
            weight="bold"
          >
            {model}
          </Text>
        </ModelChange>
      </Col>
      <Spacer y={4} />
      {/* @ts-expect-error */}
      <Card gap={2}>
        <Card.Body>
          {/* @ts-expect-error */}
          <Col justify="center" align="center">
            <Text css={{
              color: '$accents7'
            }}>请输入 Twitter ID [如 https://twitter.com/jack 即应输入 jack]</Text>
            <Spacer y={0.5}/>
            <Input
              clearable
              id={useId()}
              contentRightStyling={ false }
              aria-label='twitter id input'
              placeholder="L_x_x_x_x_x"
              onChange={handleInputChange}
              value={input}
              contentRight={
                !isLoading ? (
                  input === twitterId ? (
                    <SendButton onClick={handleRetry}>
                      <GrPowerCycle />
                    </SendButton>
                  ) : (
                    <SendButton onClick={handler}>
                      <SendIcon />
                    </SendButton>
                  )
                ) : (
                  <Loading size="sm" css={{ margin: '.5em' }} />
                )
              }
            />
            <div className="btn join-item">
              {!isLoading ? (
                <SendButton onClick={handler} className="border-none">
                  <SendIcon />
                </SendButton>
              ) : (
                <span className="loading loading-spinner"></span>
              )}
            </div>
          </div>
          {isLoading && (
            <span className="mt-4">
              由于当前请求人数较多，可能需要等待一段时间~
            </span>
          )}
          <div className="justify-center item-center">
            <div className="m-2 text-black font-semibold">
              <Suspense
                fallback={<span className="loading loading-spinner"></span>}
              >
                <Content />
              </Suspense>
            </div>
          </div>
        </div>
      </div>

      <dialog className="modal" ref={modalRef}>
        <form method="dialog" className="modal-box">
          <h3 className="font-bold text-lg">
            QaQ 可以请你帮我点一个&nbsp; <span>关注</span> &nbsp;耶
          </h3>
          <div className="py-4"></div>
          <span>你只需要点一下</span>
          <Link
            href="https://twitter.com/L_x_x_x_x_x"
            target="_blank"
            className="inline-block link link-info"
          >
            &nbsp;这里&nbsp;
          </Link>
          <span>就好啦，非常感谢~</span>
          <div className="modal-action">
            <button className="btn" onClick={handleConfirm}>
              点这里继续～
            </button>
          </div>
        </form>
      </dialog>
    </div>
  );
}
