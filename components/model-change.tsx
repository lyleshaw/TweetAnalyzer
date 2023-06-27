
import { useCallback, useState } from "react";
export type Model = {
  [x: string]: string;
};
export type ModeChangeProps = {
    models: {
        [x:string]: string
    };
    defaultSelect?: string;
    onSelect?: (mode: Model) => void;
    children?: any;
}
export function ModelChange(props: ModeChangeProps){
    const {models: modes} = props;
    const [visible, setVisible] = useState(false);
    const [selectedName, setSelectedName] = useState(props.defaultSelect);
    const handleClick = useCallback((event: React.MouseEvent<HTMLDivElement, MouseEvent> | undefined)=>{
        setVisible(!visible)
    }, [visible])
    const handleSelect = useCallback((name: string) => {
        setSelectedName(name);
        props.onSelect?.({
            [name]: modes[name]
        })
    }, [modes, props])
    return (
        <div onClick={handleClick} className={'inline-block cursor-pointer relative'}>
            {props.children}
            {
                visible && (
                    <ul className={`
                        w-full max-w-xs h-auto bg-white rounded-lg text-left p-4
                        mb-0 mt-2 mx-auto border-gray-200 border border-solid
                        absolute z-10 shadow-sm
                        dark:bg-gray-900 dark:border-0
                    `}>
                        {
                            Object.entries(modes).map(([name,api]) => {
                                return (
                                    <li key={name} className={
                                        `
                                        p-3 m-0 list-none box-border rounded-md cursor-pointer transition-all
                                        mb-2 hover:bg-gray-200 text-gray-700 ${name === selectedName ? 'bg-gray-200 text-gray-950 font-semibold' : ''}
                                        dark:hover:bg-gray-800 dark:text-gray-200 ${name === selectedName ? 'dark:bg-gray-800 dark:text-gray-100' : ''}
                                        `
                                    } onClick={()=>handleSelect(name)}>
                                        {name}
                                    </li>
                                )
                            })
                        }
                    </ul>
                )
            }
        </div>
    )
}
