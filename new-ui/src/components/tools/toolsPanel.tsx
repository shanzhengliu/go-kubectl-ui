import { Card } from "flowbite-react";
import { DockerShell } from "./dockerShell";
import { useState } from "react";
import { JsonFormatter } from "./jsonFormatter";
import { HttpHelper } from "./httpHelper";
export const ToolsPanel = () => {
  const tooksMap: { [key: string]: any } = {
    dockerShell: 
          <DockerShell iframeKey={Date.now()} />,
    jsonFormatter:
         <JsonFormatter />,
    httpHelper:
            <HttpHelper />,  
    default:null,          
  };

  const [currentTool, setCurrentTool] = useState("default");

  return (
    <div>
      <div className="flex justify-center">
        <div className="grid grid-cols-4 gap-4 m-4 max-w-screen-lg mx-auto">
          <div className=" h-20 flex justify-center items-center">
            <Card className="max-w-sm">
              <a
                onClick={() => {
                  setCurrentTool("dockerShell");
                }}
              >
                <h5 className="text-2xl font-bold tracking-tight text-gray-900">
                  Docker Shell
                </h5>
              </a>
            </Card>
          </div>
          <div className=" h-20 flex justify-center items-center">
            <Card className="max-w-sm">
              <a
                onClick={() => {
                  setCurrentTool("jsonFormatter");
                }}
              >
                <h5 className="text-2xl font-bold tracking-tight text-gray-900 ">
                  Json Formatter
                </h5>
              </a>
            </Card>
          </div>
          <div className=" h-20 flex justify-center items-center">
            <Card className="max-w-sm">
              <a
                onClick={() => {
                  setCurrentTool("httpHelper");
                }}
              >
                <h5 className="text-2xl font-bold tracking-tight text-gray-900">
                  Http Helper
                </h5>
              </a>
            </Card>
          </div>
        </div>
      </div>
      <div >{tooksMap[currentTool] ? tooksMap[currentTool] : null}</div>
    </div>
  );
};
