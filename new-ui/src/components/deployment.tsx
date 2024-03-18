import { useEffect, useState } from "react";
import { axiosInstance } from "../utils/axios";
import { DEPLOYMENT, DEPLOYMENTYAML } from "../utils/endpoints";
import { DisplayTable } from "./displayTable";
import { Button } from "flowbite-react";

export function Deployment() {
  const [renderData, setRenderData] = useState<any[]>([]);
  const fetchData = () => {
    axiosInstance
      .get(DEPLOYMENT, {
        data: {},
        headers: {
          "Content-Type": "application/json",
        },
      })
      .then((response) => {
        const responseData: any[] = [];
        for (let i = 0; i < response.data.length; i++) {
          for (let j = 0; j < response.data[i].containers.length; j++) {
            responseData.push([
              response.data[i].name,
              response.data[i].containers[j].name,
              response.data[i].containers[j].image,
              response.data[i].selector,
              response.data[i].status == 1 ? "Running" : "Pending",
              <Button>
                <a
                  href={DEPLOYMENTYAML + "?deployment=" + response.data[i].name}
                  target="_blank"
                >
                  Yaml
                </a>
              </Button>,
            ]);
          }
          setRenderData(responseData);
        }
      });
  };

  useEffect(() => {
    fetchData();
  }, []);

  return (
    <div>
      <DisplayTable
        header={["Deployment", "Container", "Image", "Selector", "Status", ""]}
        data={renderData}
        refresh={fetchData}
      ></DisplayTable>
    </div>
  );
}
