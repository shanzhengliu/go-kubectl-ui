import { useEffect, useState } from "react";
import { DisplayTable } from "./displayTable";
import { axiosInstance } from "../utils/axios";
import { POD, PODLOGS, PODYAML, WEBSHELL } from "../utils/endpoints";
import { Button } from "flowbite-react";

export function Pod() {
    const [tableData, setTableData] = useState<any[][]>([]);
    useEffect(() => {
       axiosInstance.get(POD, {
          data: {},
          headers: {
            "Content-Type": "application/json",
          }
        }).then((response) => {
          const responseData: any[][] = []
          for (let i = 0; i < response.data.length; i++) {
             for(let j=0;j<response.data[i].images.length;j++){
              responseData.push([
                response.data[i].name,
                response.data[i].namespace,
                response.data[i].images[j].containerName,
                response.data[i].images[j].name,
                <ImageId  id={response.data[i].images[j].id} ></ImageId>,
                response.data[i].images[j].containerStatus,
                response.data[i].createTime,
                <Logs pod={response.data[i].name} container={response.data[i].images[j].containerName} />,
                 <Yaml pod={response.data[i].name} />,
                <WebShell pod={response.data[i].name} container={response.data[i].images[j].containerName} />]);
            }
          }
          setTableData(responseData);
        });
      }, []);
   
    return <div>
      <DisplayTable data={tableData} header={["Pod", "Namespace", "Container","Image","Image Id","Status","Create Time","","",""]} />
    </div>;
  }

  
function Logs(props: {pod: string, container: string}) {
  return (
    <Button gradientMonochrome="info">
    <div>
       <a href={`${PODLOGS}?pod=${props.pod}&container=${props.container}`} target="_blank">Logs</a>
    </div>
    </Button>
  );
}

function WebShell(props: {pod: string, container: string}) {
  return (
    <div>
      <Button gradientMonochrome="success">
       <a href={`${WEBSHELL}?pod=${props.pod}&container=${props.container}`} target="_blank">Shell</a>
       </Button>
    </div>
  );
}

function ImageId(props: {id: string}) {
  return (
    <div>
            <Button color="light" onClick={() => alert(props.id)}>ID</Button>
    </div>
  );
}

function Yaml(props: {pod: string}) {
  return (
    <div>
      <Button gradientDuoTone="cyanToBlue">
       <a href={`${PODYAML}?pod=${props.pod}`} target="_blank">Yaml</a>
       </Button>
    </div>
  );
}