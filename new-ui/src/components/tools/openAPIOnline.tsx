import {Button, FileInput, Label, TextInput} from "flowbite-react";
import {useEffect, useState} from "react";
import {
    OPENAPI_HELPER_LIST,
    OPENAPI_HELPER_START,
    OPENAPI_HELPER_STOP,
    OPENAPI_HELPER_STOP_ALL,
    OPENAPI_HELPER_UPLOAD
} from "../../utils/endpoints";

import {FileTreeComponent} from "./fileTree";
import {inputHook} from "../../hooks/inputhook";
import {axiosInstance} from "../../utils/axios";
import Swal from "sweetalert2";
import {errorHander} from "../../utils/tools.tsx";

interface OpenAPIList {
    path: string;
    port: string;
    
}
export const OpenAPIOnline = () => {
  const [file, setFile] = useState(null);
  const [selectFile, setSelectFile] = useState<string>("NONE");
  const [uploadTimestamp, setUploadTimestamp] = useState(Date.now());
  const [port, , onChangeSetPort] = inputHook("7001");
  const [openApiList, setOpenApiList] = useState<OpenAPIList[]>([]);

  const selectFileSelectForListener = (path: string) => {
    setSelectFile(path);
  };

  const handleFileChange = (event: any) => {
    setFile(event.target.files[0]);
  };

  const uploadZip = async (file: any) => {
    const formData = new FormData();
    formData.append("file", file);
    try {
      await fetch(OPENAPI_HELPER_UPLOAD, {
        method: "POST",
        body: formData,
      });
      setUploadTimestamp(Date.now());
    } catch (error) {
        errorHander(error)
    }
  };

  const startListener = async () => {
      if (selectFile === "NONE") {
          await Swal.fire({
                icon: "error",
                text: "Please select a file to start listener",
            });
          return ;
          }

     await listenOperation("start", selectFile, port);
     setOpenApiList( await axiosInstance({ method: "GET", url: OPENAPI_HELPER_LIST }).then(response => response.data)) ;
  }

  useEffect(() => {

      axiosInstance({ method: "GET", url: OPENAPI_HELPER_LIST })
        .then(response => {
            if (response.data !== null) {
                setOpenApiList(response.data);
            }
          
        })
    
    
    
  }, []);

  const stopAllListener = async () => {
    await axiosInstance({ method: "GET", url: OPENAPI_HELPER_STOP_ALL });
    setOpenApiList( await axiosInstance({ method: "GET", url: OPENAPI_HELPER_LIST }).then(response => response.data)) ;
  }

  const stopListener = async (filename:string, port:string) => {
    await listenOperation("stop", filename, port);
    setOpenApiList( await axiosInstance({ method: "GET", url: OPENAPI_HELPER_LIST }).then(response => response.data)) ;
  }

  const listenOperation = async (operation:string, filename:string, port: string) => {
    if (filename == "NONE") {
        return;
    }
    const url = operation === "start" ? OPENAPI_HELPER_START : OPENAPI_HELPER_STOP;
    const options = {
        method: "POST",
        url: url,
        headers: {
            "Content-Type": "application/json",
        },
        data: JSON.stringify({
            port: port,
            path: "/tmp/kubectl-go-upload"+filename,
        }),
    };
    try {
        return await axiosInstance(options);
    } catch (error) {
        errorHander(error)
        return error;
    } 
  }
  

  return (
    <div className="">
        <div>
            <h1>OpenAPI Online </h1>
            <span>you can upload the openapi  <span style={{color:"red"}}>zip</span>, <span style={{color:"red"}}>yaml</span> file and start a server for testing via ui, now you can set up <span style={{"color":"red"}}>"openapi-status-code"</span>, <span style={{"color":"red"}}>"openapi-example"</span>, or <span style={{"color":"red"}}>"openapi-content-type"</span> in request header to control the response </span>
        </div>
      <div className="mt-4 mb-2 w-full flex">
        <div className="mt-2 block">
          Upload file
        </div>
        <div className="flex-grow ml-2">
        <FileInput
          accept=".zip,.yaml,.yml"
          className="flex-grow"
          id="file-upload"
          onChange={handleFileChange}
        />
        </div>
        <div>
        <Button className="ml-2" onClick={() => uploadZip(file)}>
          Submit
        </Button>
        </div>
      </div>
      <div className="mb-2 w-full flex flex-grow">
        <div className="mb-2 w-full">
          <FileTreeComponent
            key={uploadTimestamp}
            setSelectFile={selectFileSelectForListener}
          />
        </div>
        <div className="w-full ml-2 mt-2">
          <Label value="Selected File" />
          <div className="p-2 mt-2">{selectFile}</div>
          <Label value="Listen Port" />
          <TextInput
            className="w-full mt-2"
            value={port}
            onChange={onChangeSetPort}
          />
          <div className="flex mt-2 justify-end" >
          <Button className="mt-4" onClick={startListener} >Start</Button>
          </div>
         
        </div>
        <div className="w-full ml-2 mt-2">
           <div className="flex justify-between">
          <h2>Running OpenAPI Mock</h2>
          <Button className="mt-4" color={"failure"} onClick={stopAllListener}  >Stop All</Button>
          </div>
            <div className="p-2">
                {openApiList && openApiList.map((item, index) => {
                return (
                    <div key={index} className="flex justify-between">
                    <div>{item.path} - {item.port}</div>
                    <Button  color="warning" className="mt-4" onClick={() => stopListener(item.path, item.port)}>Stop</Button>
                    </div>
                );
                })}
          </div>
        </div>
      </div>
    </div>
  );
};
