import { useState, ChangeEvent } from 'react';
export const inputHook = (initialValue: string) => {
    const [value, setValue] = useState(initialValue);
    const onChange = (e: ChangeEvent<HTMLInputElement>) => {
        setValue(e.target.value);
    };
    return [value, setValue, onChange] as [string, (e: string) => void, (e: ChangeEvent<HTMLInputElement>) => void];
}


    