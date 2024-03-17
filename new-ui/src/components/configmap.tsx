import { useEffect, useState } from "react";
import { DisplayTable } from "./displayTable";
import { axiosInstance } from "../utils/axios";
import { CONFIGMAP } from "../utils/endpoints";
export function Configmap() {
  const [tableData, setTableData] = useState<any[][]>([]);

  useEffect(() => {
    axiosInstance
      .get(CONFIGMAP, {
        data: {},
        headers: {
          Authorization: "hello",
          "Content-Type": "application/json",
        },
      })
      .then((response) => {
        const responseData: any[][] = []
        response.data.map((item: any) => {
          responseData.push([item.name, item.namespace]);
        
        });
        setTableData(responseData);
      });
  }, []);
  return ( 
    <DisplayTable data={tableData} header={["Configmap", "Namespace", ""]} />
  );
}
