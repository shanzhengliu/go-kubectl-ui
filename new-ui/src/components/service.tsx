import { useEffect, useState } from "react";
import { DisplayTable } from "./displayTable";
import { axiosInstance } from "../utils/axios";
import { SERVICE } from "../utils/endpoints";

export function Service() {
  const [renderData, setRenderData] = useState<any[][]>([]);
  const fetchData = () => {
    axiosInstance
    .get(SERVICE, {
      data: {},
      headers: {
        "Content-Type": "application/json",
      },
    })
    .then((response) => {
      const responseData: any[] = [];
      for (let i = 0; i < response.data.length; i++) {
        responseData.push([
          response.data[i].name,
          response.data[i].namespace,
          response.data[i].type,
          response.data[i].selector,
        
        ]);
      }
      setRenderData(responseData);
    });
  }
   
   useEffect(() => {
    fetchData();
   }, []);

    const refresh = () => {
      fetchData();
    }


  return (
    <div>
      <DisplayTable
        header={["Service", "NameSpace", "Type", "Selector"]}
        data={renderData}
        refresh={refresh}
      ></DisplayTable>
    </div>
  );
}
