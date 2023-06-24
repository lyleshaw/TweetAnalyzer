import styled from "@emotion/styled";
import React, {
  Suspense,
  useCallback,
  useState,
  useTransition,
  useRef,
} from "react";
import { atomWithObservable } from "jotai/utils";
import { useAtomValue, useSetAtom } from "jotai/react";
import { atom } from "jotai/vanilla";
import { Observable } from "rxjs";
import Link from "next/link";

export const SendButton = styled.button`
  background: "transparent";
  border: none;
  padding: 0;

  width: 20px;
  margin: "0 10px";
  display: flex;
  justify-content: center;

  border-radius: 50%;
  cursor: pointer;
  transition: opacity 0.25s ease 0s, transform 0.25s ease 0s;
  svg {
    size: "100%";
    padding: 4px;
    transition: transform 0.25s ease 0s, opacity 200ms ease-in-out 50ms;
    box-shadow: 0 5px 20px -5px rgba(0, 0, 0, 0.1);
  }
  &:hover {
    opacity: 0.8;
  }
  &:active {
    transform: scale(0.9);
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

const contentAtom = atomWithObservable((get) => {
  const id = get(twitterIdAtom);
  return new Observable<string | null>((subscriber) => {
    const abortController = new AbortController();
    if (id === null) {
      // no value
      subscriber.next(null);
    } else {
      async function fetchData() {
        const response = await fetch(
          "https://tweet-api.aireview.tech/api/get_tweet_analysis?twitter_id=" +
            id,
          {
            headers: {
              "Content-Type": "application/json",
            },
            signal: abortController.signal,
            method: "GET",
          }
        );
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

      fetchData().catch(subscriber.error);
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
  const [input, setInput] = useState("");
  const modalRef = useRef<HTMLDialogElement>(null);

  const [isLoading, startTransition] = useTransition();
  const setTwitterId = useSetAtom(twitterIdAtom);
  const handleInputChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setInput(event.target.value);
    },
    []
  );

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

  return (
    <div className="lg:m-12 sm:m-8 m-8">
      <div className="text-center mt-24">
        <div className="indicator">
          <span className="indicator-item badge -right-4 -top-4 ">
            OpenAI 版
          </span>
          <div className="inline bg-gradient-to-r from-[rgba(0,110,58,0.8)]  to-[rgba(22,107,181)] bg-clip-text font-display text-5xl tracking-tight text-transparent">
            赛博算命师
          </div>
        </div>
        <div className="block bg-gradient-to-r from-[rgba(0,110,58,0.8)]  to-[rgba(22,107,181)] bg-clip-text font-display text-5xl tracking-tight text-transparent">
          Cyber Fortune Teller
        </div>
      </div>

      <div className="card shadow-xl bg-base-100 w-auto">
        <div className="card-body text-center p-2 py-12">
          <div className="text-gray-500">
            请输入 Twitter ID
            <div>[如 https://twitter.com/jack 即应输入 jack]</div>
          </div>
          <div className="join mx-auto">
            <button className="btn join-item">ID</button>
            <input
              className="input input-bordered join-item "
              placeholder="L_x_x_x_x_x"
              onChange={handleInputChange}
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
