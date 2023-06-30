import styled from "@emotion/styled";
import React, {
  Suspense,
  useCallback,
  useState,
  useTransition,
  useRef,
  type ComponentPropsWithoutRef,
} from "react";
import { atomWithObservable, atomWithStorage } from "jotai/utils";
import { useAtom, useAtomValue, useSetAtom } from "jotai/react";
import { atom } from "jotai/vanilla";
import { Observable } from "rxjs";
import Link from "next/link";
import { Model, ModelChange } from "../components/model-change";

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

// GrPowerCycle from react-icons
const RetryIcon = (props: ComponentPropsWithoutRef<"svg"> = {}) => (
  <svg
    stroke="currentColor"
    fill="currentColor"
    strokeWidth="0"
    viewBox="0 0 24 24"
    height="24"
    width="24"
    xmlns="http://www.w3.org/2000/svg"
    {...props}
  >
    <path
      fill="none"
      strokeWidth="2"
      d="M20,8 C18.5343681,5.03213345 15.4860999,3 11.9637942,3 C7.01333514,3 3,7.02954545 3,12 M4,16 C5.4656319,18.9678666 8.51390007,21 12.0362058,21 C16.9866649,21 21,16.9704545 21,12 M9,16 L3,16 L3,22 M21,2 L21,8 L15,8"
    ></path>
  </svg>
);

const twitterIdAtom = atom<string | null>(null);
const models: Record<string, string> = {
  "Open-AI": "https://tweet-api-boe.aireview.tech/api/get_tweet_analysis",
  Claude: "https://tweet-api.aireview.tech/api/get_tweet_analysis",
};
const modelAtom = atom<string>("Claude");

const contentQueryRequestIdStateAtom = atom<number>(0);

const contentAtom = atomWithObservable((get) => {
  const id = get(twitterIdAtom);
  const modelName = get(modelAtom);
  const api = models[modelName];
  get(contentQueryRequestIdStateAtom);
  return new Observable<string | null>((subscriber) => {
    const abortController = new AbortController();
    if (id === null) {
      // no value
      subscriber.next(null);
    } else {
      async function fetchData() {
        const response = await fetch(`${api}?twitter_id=${id}`, {
          headers: {
            "Content-Type": "application/json",
          },
          signal: abortController.signal,
          method: "GET",
        });
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

      fetchData().catch((err) => subscriber.error(err));
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
  const [model, setModel] = useAtom(modelAtom);
  const onSelect = (model: Model) => {
    const name = Object.keys(model)[0];
    setModel(name);
  };
  const [requestId, setRequestId] = useAtom(contentQueryRequestIdStateAtom);

  const handleInputChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setInput(event.target.value);
    },
    []
  );

  const handleConfirm = useCallback(() => {
    startTransition(() => {
      if (input !== "") {
        setTwitterId(input);
      }
      setRequestId((id) => (id += 1));
    });
  }, [input, setRequestId, setTwitterId]);

  const handleClick = useCallback(() => {
    if (!requestId) {
      modalRef.current?.showModal();
      return;
    }

    handleConfirm();
  }, [handleConfirm, requestId]);

  return (
    <div className="m-8 lg:m-12 sm:m-8">
      <div className="mt-24 mb-8 text-center">
        <ModelChange
          onSelect={onSelect}
          defaultSelect={model}
          models={models}
        >
          <div className="inline flex content-center">
            <div className="flex items-center">
              <svg 
                className="bi bi-caret-down-fill bg-gradient-to-r from-[rgba(0,110,58,0.8)]  to-[rgba(22,107,181)] bg-clip-text" 
                xmlns="http://www.w3.org/2000/svg" 
                width="16" 
                height="16" 
                fill="currentColor" 
                viewBox="0 0 16 16">
                <defs>
                  <radialGradient  id="trangle">
                      <stop offset="0%" stop-color="#006e3a" />
                      <stop offset="100%" stop-color="#166bb5" />
                  </radialGradient>
                </defs>
                <path 
                  d="M7.247 11.14 2.451 5.658C1.885 5.013 2.345 4 3.204 4h9.592a1 1 0 0 1 .753 1.659l-4.796 5.48a1 1 0 0 1-1.506 0z"
                  stroke="none"
                  fill="url(#trangle)"/>
              </svg>
              <div className="bg-gradient-to-r from-[rgba(0,110,58,0.8)]  to-[rgba(22,107,181)] bg-clip-text font-display text-5xl tracking-tight text-transparent">
                {model}
              </div>
            </div>
            <p className="bg-gradient-to-r from-[rgba(0,110,58,0.8)]  to-[rgba(22,107,181)] bg-clip-text font-display text-sm tracking-tight text-transparent">
              版
            </p>
          </div>
        </ModelChange>
        <div className="block bg-gradient-to-r from-[rgba(0,110,58,0.8)]  to-[rgba(22,107,181)] bg-clip-text font-display text-5xl tracking-tight text-transparent">
          赛博算命师
        </div>
        <div className="block bg-gradient-to-r from-[rgba(0,110,58,0.8)]  to-[rgba(22,107,181)] bg-clip-text font-display text-5xl tracking-tight text-transparent">
          Cyber Fortune Teller
        </div>
      </div>

      <div className="w-auto text-gray-500 shadow-xl card bg-base-100 dark:bg-slate-800 dark:text-slate-50">
        <div className="p-2 py-12 text-center card-body">
          <div>
            请输入 Twitter ID
            <div>[如 https://twitter.com/jack 即应输入 jack]</div>
          </div>
          <div className="mx-auto join">
            <button className="btn join-item">ID</button>
            <input
              className="input input-bordered join-item"
              placeholder="L_x_x_x_x_x"
              style={{ width: "100%" }}
              onChange={handleInputChange}
            />
            <div className="btn join-item">
              {!isLoading ? (
                <SendButton onClick={handleClick} className="border-none">
                  {requestId ? <RetryIcon /> : <SendIcon />}
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
            <div className="m-2 font-semibold">
              <Suspense
                fallback={<span className="loading loading-spinner"></span>}
              >
                <Content />
              </Suspense>
            </div>
          </div>
        </div>
      </div>
      <footer className="p-10 footer text-neutral-content">
        <div>
          <p>Cyber Fortune Teller</p>
        </div>
        <div>
          <div className="grid grid-flow-col gap-4">
            <Link href="https://twitter.com/L_x_x_x_x_x">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                className="fill-current"
              >
                <path d="M24 4.557c-.883.392-1.832.656-2.828.775 1.017-.609 1.798-1.574 2.165-2.724-.951.564-2.005.974-3.127 1.195-.897-.957-2.178-1.555-3.594-1.555-3.179 0-5.515 2.966-4.797 6.045-4.091-.205-7.719-2.165-10.148-5.144-1.29 2.213-.669 5.108 1.523 6.574-.806-.026-1.566-.247-2.229-.616-.054 2.281 1.581 4.415 3.949 4.89-.693.188-1.452.232-2.224.084.626 1.956 2.444 3.379 4.6 3.419-2.07 1.623-4.678 2.348-7.29 2.04 2.179 1.397 4.768 2.212 7.548 2.212 9.142 0 14.307-7.721 13.995-14.646.962-.695 1.797-1.562 2.457-2.549z"></path>
              </svg>
            </Link>
            <Link href="https://github.com/lyleshaw/TweetAnalyzer">
              <svg
                viewBox="0 0 16 16"
                className="w-5 h-5"
                fill="currentColor"
                aria-hidden="true"
              >
                <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z" />
              </svg>
            </Link>
          </div>
        </div>
      </footer>

      <dialog className="modal" ref={modalRef}>
        <form method="dialog" className="modal-box">
          <h3 className="text-lg font-bold">
            QaQ 可以请你帮我点一个&nbsp; <span>关注</span> &nbsp;耶
          </h3>
          <div className="py-4"></div>
          <div className="flex">
            你只需要点一下
            <Link
              href="https://twitter.com/L_x_x_x_x_x"
              target="_blank"
              className="inline-block px-2 link link-info"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                className="fill-current"
              >
                <path d="M24 4.557c-.883.392-1.832.656-2.828.775 1.017-.609 1.798-1.574 2.165-2.724-.951.564-2.005.974-3.127 1.195-.897-.957-2.178-1.555-3.594-1.555-3.179 0-5.515 2.966-4.797 6.045-4.091-.205-7.719-2.165-10.148-5.144-1.29 2.213-.669 5.108 1.523 6.574-.806-.026-1.566-.247-2.229-.616-.054 2.281 1.581 4.415 3.949 4.89-.693.188-1.452.232-2.224.084.626 1.956 2.444 3.379 4.6 3.419-2.07 1.623-4.678 2.348-7.29 2.04 2.179 1.397 4.768 2.212 7.548 2.212 9.142 0 14.307-7.721 13.995-14.646.962-.695 1.797-1.562 2.457-2.549z"></path>
              </svg>
            </Link>
            <span>就好啦，非常感谢~</span>
          </div>
          <div className="flex">
            如果可以的话，可以再给我一个 🌟
            <Link href="https://github.com/lyleshaw/TweetAnalyzer">
              <svg
                viewBox="0 0 16 16"
                className="w-5 h-5"
                fill="currentColor"
                aria-hidden="true"
              >
                <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z" />
              </svg>
            </Link>
            <span>吗？~~</span>
          </div>
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
