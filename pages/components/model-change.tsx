import styled from "@emotion/styled";
import { useCallback, useState } from "react";
export type Model = {
  [x: string]: string;
};
export type ModeChangeProps = {
  models: {
    [x: string]: string;
  };
  defaultSelect?: string;
  onSelect?: (mode: Model) => void;
  children?: any;
};
export const ModelChangeWrapper = styled.div`
  display: inline-block;
  cursor: pointer;
`;

export const ModelChangeMenu = styled.ul`
  width: 250px;
  height: auto;
  background-color: #fff;
  borderradius: "$md";
  text-align: left;
  padding: 4px;
  margin: 0 auto;
  border: "$accents2 1px solid";
`;
const Model = styled.li`
  padding: "$3 $3";
  margin: "0";
  list-style: none;
  box-sizing: border-box;
  border-radius: "$md";
  cursor: "pointer";
  margin-bottom: "$2";
  &[data-active="true"] {
    background: "$gray100";
  }
  &:hover {
    background: "$gray100";
  }
`;

export function ModelChange(props: ModeChangeProps) {
  const { models: modes } = props;
  const [visible, setVisible] = useState(false);
  const [selectedName, setSelectedName] = useState(props.defaultSelect);
  const handleClick = useCallback(
    (event: React.MouseEvent<HTMLDivElement, MouseEvent> | undefined) => {
      setVisible(!visible);
    },
    [visible]
  );
  const handleSelect = useCallback(
    (name: string) => {
      setSelectedName(name);
      props.onSelect?.({
        [name]: modes[name],
      });
    },
    [setSelectedName]
  );
  return (
    <ModelChangeWrapper onClick={handleClick}>
      {props.children}
      {visible && (
        <ModelChangeMenu>
          {Object.entries(modes).map(([name, api]) => {
            return (
              <Model
                key={name}
                onClick={() => handleSelect(name)}
                data-active={name === selectedName}
              >
                {name}
              </Model>
            );
          })}
        </ModelChangeMenu>
      )}
    </ModelChangeWrapper>
  );
}
