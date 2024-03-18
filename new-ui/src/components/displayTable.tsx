import {
  Button,
  Checkbox,
  Label,
  Modal,
  Select,
  Table,
  TextInput,
} from "flowbite-react";
import React, { useEffect } from "react";
import {
  JSXElementConstructor,
  Key,
  ReactElement,
  ReactNode,
  ReactPortal,
} from "react";
import { axiosInstance } from "../utils/axios";
import { CONTEXT_CHANGE, CONTEXT_LIST, CURRENT_CONTEXT } from "../utils/endpoints";
import { inputHook } from "../hooks/inputhook";
import { ContextSwitcher } from "./contextSwitcher";

export function DisplayTable(props: {
  header: any[];
  data: any[][];
  checkbox?: boolean;
  refresh: () => void;
}) {
  const [searchValue, setSearchValue] = React.useState("");
  const [currentContext, setCurrentContext] = React.useState("NA");
  const [contextSelected, setContextSelected] = React.useState("");
  const [currentNamespace, setCurrentNamespace] = React.useState("NA");
  const [inputNamespace, setInputNamespace, onchangeInputNamespace] = inputHook("");
  const [renderData, setRenderData] = React.useState<any[][]>(props.data || []);
  const [originData, setOriginData] = React.useState<any[][]>(props.data || []);
  const [contextList, setContextList] = React.useState<any[]>([]);
  const [openModal, setOpenModal] = React.useState(false);
  const onSearchChange = (e: {
    target: { value: React.SetStateAction<string> };
  }) => {
    setSearchValue(e.target.value);
    if (e.target.value === "") {
      setRenderData(originData);
    } else {
      const result = originData.filter((row) => {
        return row.some((cell) => {
          cell = cell.toString();
          return cell.includes(searchValue);
        });
      });
      setRenderData(result);
    }
  };

  const contextValueChange = (e:any) => {
    setContextSelected(e.target.value);
  }


  useEffect(() => {
    if (openModal) {
    setContextSelected(currentContext);
    }
    }, [openModal]);

  useEffect(() => {
    if (searchValue === "") {
      setRenderData(originData);
    } else {
      const result = originData.filter((row) => {
        return row.some((cell) => {
          let cellString = "";
          if (typeof cell !== "string" || typeof cell === "object") {
            cell = JSON.stringify(cell);
          }
          cellString = cell.toString().toLowerCase();

          return cellString.includes(searchValue.toLowerCase());
        });
      });
      setRenderData(result);
    }
  }, [searchValue]);

  useEffect(() => {
    setRenderData(props.data);
    setOriginData(props.data);
  }, [props.data]);

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

  const switchContext = () => {
    axiosInstance
      .get(CONTEXT_CHANGE+"?context="+contextSelected+"&namespace="+inputNamespace, {
        data: {},
        headers: {
          "Content-Type": "application/json",
        },
      })
      .then(async () => {
        setCurrentContext(contextSelected);
        setCurrentNamespace(inputNamespace);
        props.refresh();
        setOpenModal(false);
      });

   
  }

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
        <div className="flex items-center mb-4 ">
          <Button
            className="mr-2"
            color="success"
            onClick={() => {
              props.refresh && props.refresh();
            }}
          >
            Refresh
          </Button>
          <TextInput
            className="mr-2"
            placeholder="Search..."
            value={searchValue}
            onChange={onSearchChange}
          />
        </div>
      </div>
      <div className="overflow-x-auto">
        <Table hoverable>
          <Table.Head>
            {props.checkbox ? (
              <Table.HeadCell className="p-4">
                <Checkbox />
              </Table.HeadCell>
            ) : null}
            {props.header.map((header, index) => (
              <Table.HeadCell key={index}>{header}</Table.HeadCell>
            ))}
          </Table.Head>
          <Table.Body className="divide-y">
            {renderData.map((row, index) => (
              <Table.Row
                key={index}
                className="bg-white dark:border-gray-700 dark:bg-gray-800"
              >
                {props.checkbox ? (
                  <Table.HeadCell className="p-4">
                    <Checkbox />
                  </Table.HeadCell>
                ) : null}
                {row.map(
                  (
                    cell:
                      | string
                      | number
                      | boolean
                      | ReactElement<any, string | JSXElementConstructor<any>>
                      | Iterable<ReactNode>
                      | ReactPortal
                      | null
                      | undefined,
                    index: Key | null | undefined
                  ) => (
                    <Table.Cell key={index}>
                      {React.isValidElement(cell) ? cell : <span>{cell}</span>}
                    </Table.Cell>
                  )
                )}
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      </div>

      <ContextSwitcher onSwitch={props.refresh} />
    </div>
  );
}
