import { useEffect, useState } from "react";
import { axiosInstance } from "../utils/axios";
import {
  CONTEXT_CHANGE,
  CONTEXT_LIST,
  CURRENT_CONTEXT,
  RESOURCE,
} from "../utils/endpoints";
import ReactECharts from "echarts-for-react";
import { colors } from "../utils/constant";
import { Button, Label, Modal, Select, TextInput } from "flowbite-react";
export function Resource() {
  const [renderData, setRenderData] = useState<any[][]>([]);
  const [currentContext, setCurrentContext] = useState("NA");
  const [contextSelected, setContextSelected] = useState("");
  const [currentNamespace, setCurrentNamespace] = useState("NA");
  const [inputNamespace, setInputNamespace] = useState("");
  const [contextList, setContextList] = useState<any[]>([]);
  const [openModal, setOpenModal] = useState(false);
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

  const contextValueChange = (e: any) => {
    setContextSelected(e.target.value);
  };

  const inputNamespaceChange = (e: any) => {
    setInputNamespace(e.target.value);
  };

  useEffect(() => {
    axiosInstance
      .get(CURRENT_CONTEXT, {
        data: {},
        headers: {
          "Content-Type": "application/json",
        },
      })
      .then((response) => {
        setCurrentContext(response.data.context);
        setCurrentNamespace(response.data.namespace);
        setInputNamespace(response.data.namespace);
      });
    axiosInstance
      .get(CONTEXT_LIST, {
        data: {},
        headers: {
          "Content-Type": "application/json",
        },
      })
      .then((response) => {
        setContextList(response.data);
      });
  }, []);

  useEffect(() => {
    if (openModal) {
      setContextSelected(currentContext);
    }
  }, [openModal]);

  useEffect(() => {
    fetchData();
  }, []);

  const switchContext = () => {
    axiosInstance
      .get(
        CONTEXT_CHANGE +
          "?context=" +
          contextSelected +
          "&namespace=" +
          inputNamespace,
        {
          data: {},
          headers: {
            "Content-Type": "application/json",
          },
        }
      )
      .then(async () => {
        setCurrentContext(contextSelected);
        setCurrentNamespace(inputNamespace);
        fetchData();
        setOpenModal(false);
      });
  };

  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <div className="flex items-center mb-4">
          <span className="ml-2">namespace:</span>
          <span className="ml-2 text-green-600"> {currentContext}</span>
          <span className="ml-2">context:</span>
          <span className="ml-2 text-blue-600">{currentNamespace}</span>
          <Button
            className="ml-2"
            onClick={() => {
              setOpenModal(true);
            }}
          >
            Switch
          </Button>
        </div>
      </div>
      <div className="flex justify-end items-center mb-4">
        <Button className="mr-4" color={"success"} onClick={fetchData}>
          Refresh
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {renderData.map((option, index) => (
          <ReactECharts key={index} option={option} />
        ))}
      </div>

      <Modal
        show={openModal}
        onClose={() => {
          setOpenModal(false);
        }}
      >
        <Modal.Header>Context</Modal.Header>
        <Modal.Body>
          <div className="max-w-md">
            <div className="mb-2 block">
              <Label htmlFor="context" value="Select your context" />
            </div>
            <Select
              id="context"
              required
              value={contextSelected}
              onChange={contextValueChange}
            >
              {contextList.map((context, index) => (
                <option value={context} key={index}>
                  {context}
                </option>
              ))}
            </Select>
            <div className="mb-2 block">
              <Label htmlFor="context" value="Input your Namespace" />
              <TextInput
                id="context"
                value={inputNamespace}
                onChange={inputNamespaceChange}
                className="bg-gray-200"
                placeholder="Input your Namespace"
              />
            </div>
          </div>
        </Modal.Body>
        <Modal.Footer>
          <Button onClick={switchContext}>OK</Button>
        </Modal.Footer>
      </Modal>
    </div>
  );
}
