import { useEffect, useState } from "react";
import { axiosInstance } from "../utils/axios";
import { RESOURCE } from "../utils/endpoints";
import ReactECharts from "echarts-for-react";
import { colors } from "../utils/constant";

import { ContextSwitcher } from "./contextSwitcher";
export function Resource() {
  const [renderData, setRenderData] = useState<any[][]>([]);
  const textFirstCharCapitalize = (value: string) => {
    return value.replace(".", " ").replace(/\b\w/g, (l) => l.toUpperCase());
  };

  const shuffle = (array: string[]) => {
    for (let i = array.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [array[i], array[j]] = [array[j], array[i]];
    }
    return array;
  };

  const generateOption = (key: string, resource: { [key: string]: any }) => {
    const randomColor = shuffle(colors);
    const data = [
      { name: "used", value: resource["used"], itemStyle: { color: "red" } },
      { name: "free", value: resource["free"], itemStyle: { color: "green" } },
    ];

    data.forEach(function (item, index) {
      item.itemStyle = {
        color: randomColor[index],
      };
    });
    var option = {
      title: {
        text: textFirstCharCapitalize(key),
        left: "center",
      },
      tooltip: {
        trigger: "item",
      },
      series: [
        {
          name:
            textFirstCharCapitalize(key) +
            `Resource Usage - Total : ${resource["hard"]} ${resource["unit"]}`,
          type: "pie",
          radius: "60%",
          data: data,
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: "rgba(0, 0, 0, 0.5)",
            },
          },
          label: {
            formatter: `{b}: {c} ${resource["unit"]} {d}%`,
          },
        },
      ],
    };

    return option;
  };
  const fetchData = () => {
    axiosInstance
      .get(RESOURCE, {
        data: {},
        headers: {
          "Content-Type": "application/json",
        },
      })
      .then((response) => {
        const responseData: any[] = [];
        for (let i = 0; i < response.data.length; i++) {
          const resourceMap = response.data[i].resourceMap;
          Object.keys(resourceMap).forEach((key) => {
            responseData.push(generateOption(key, resourceMap[key]));
          });
        }

        setRenderData(responseData);
      });
  };

  useEffect(() => {
    fetchData();
  }, []);

  return (
    <div>
      <ContextSwitcher onSwitch={fetchData} />
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {renderData.map((option, index) => (
          <ReactECharts key={index} option={option} />
        ))}
      </div>
    </div>
  );
}
