import {Button, TextInput } from 'flowbite-react';
import { LOCALSHELL } from '../utils/endpoints';
export const ToolSet = (props: {
    getKeyWord: any[];
  }
   
) => {
  return (
   <div className="flex justify-end items-center mb-4 ">
      <Button className="mr-4"><a href={LOCALSHELL} target='_blank'  >Docker</a></Button>  
      <TextInput  className="mr-2" placeholder="Search..."  />
    </div>
  );
};
