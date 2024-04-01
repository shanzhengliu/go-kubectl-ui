import axios from "axios";
import { Button, Select, TextInput, Textarea } from "flowbite-react";
import { Tabs } from "flowbite-react";
import { inputHook } from "../../hooks/inputhook";
import { PROXY } from "../../utils/endpoints";
import { useState } from "react";
import { html as beautifyHtml } from "js-beautify";

type Header = {
  key: string;
  value: string;
};

export const HttpHelper = () => {
  const [method, , onChangeMethod] = inputHook("GET");
  const [requestUrl, , onChangeRequestUrl] = inputHook("");
  const [responseData, setResponseData] = inputHook("");
  const [responseStatus, setResponseStatus] = inputHook("");
  const [requestBody, , onChangeRequestBody] = inputHook("{}");
  const [responseHeader, setResponseHeader] = useState<Header[]>([]);
  const [columnsData, setColumnsData] = useState<Header[]>([
    { key: "Content-Type", value: "application/json" },
    { key: "Accept", value: "*/*" },
    {
      key: "User-Agent",
      value:
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36",
    },
  ]);
  const addColumn = () => {
    if (
      columnsData[columnsData.length - 1].key === "" ||
      columnsData[columnsData.length - 1].value === ""
    ) {
      return;
    }
    setColumnsData([...columnsData, { key: "", value: "" }]);
  };

  const deleteColumnData = (index: number) => {
    const newColumnsData = [...columnsData];
    newColumnsData.splice(index, 1);
    setColumnsData(newColumnsData);
  };

  const updateColumnData = (index: number, key: string, data: string) => {
    const newColumnsData = [...columnsData];
    if (key === "key") {
      newColumnsData[index].key = data;
    }
    if (key === "value") {
      newColumnsData[index].value = data;
    }
    setColumnsData(newColumnsData);
  };

  const requestHeaderGenerate = () => {
    const headers: { [key: string]: string } = {};
    columnsData.forEach((column) => {
      if (column.key !== "" && column.value !== "") {
        headers[column.key] = column.value;
      }
    });
    return headers;
  };

  const responseHeaderGenerate = (responseHeader: any) => {
    const headers: Header[] = [];
    for (const [key, value] of Object.entries(responseHeader)) {
      headers.push({ key: key as string, value: value as string });
    }
    return headers;
  };

  const bodyGenerate = () => {
    try {
      return JSON.parse(requestBody);
    } catch (e) {
      return requestBody;
    }
  };

  const send = () => {
    console.log(method);
    console.log(requestUrl);
    const options = {
      method: method,
      url: requestUrl,
      headers: requestHeaderGenerate(),
      body: bodyGenerate(),
    };
    axios
      .post(PROXY, options)
      .then((response) => {
        setResponseData(beautifyHtml(JSON.stringify(response.data, null, 2)));
        setResponseStatus(response.status + " " + response.statusText);
        setResponseHeader(responseHeaderGenerate(response.headers));
      })
      .catch((error) => {
        setResponseData(JSON.stringify(error.response.data, null, 2));
        setResponseStatus(
          error.response.status + " " + error.response.statusText
        );
        setResponseHeader(responseHeaderGenerate(error.response.headers));
      });

  };
  return (
    <div>
      <div className="flex-grow  flex max-w-[calc(100vw-3rem)] ">
        <div className="w-48 m-2">
          <Select value={method} onChange={onChangeMethod}>
            <option value="GET">GET</option>
            <option value="POST">POST</option>
            <option value="PUT">PUT</option>
            <option value="DELETE">DELETE</option>
            <option value="PATCH">PATCH</option>
            <option value="OPTIONS">OPTIONS</option>
          </Select>
        </div>
        <div className="w-full m-2">
          <TextInput
            type="url"
            placeholder="place your url..."
            value={requestUrl}
            onChange={onChangeRequestUrl}
          ></TextInput>
        </div>
        <div className="m-2">
          <Button color="success" onClick={send}>
            Send
          </Button>
        </div>
      </div>
      <div className="flex-grow flex w-full">
        <div id="requestBlock" className="w-full m-2">
          <Tabs>
            <Tabs.Item active title="Request Header">
              {columnsData.map((column, index) => (
                <div key={index} className="flex flex-grow p-2">
                  <TextInput
                    type="text"
                    value={column.key}
                    onChange={(e) =>
                      updateColumnData(index, "key", e.target.value)
                    }
                    className="w-full m-2"
                  />
                  <TextInput
                    type="text"
                    value={column.value}
                    onChange={(e) =>
                      updateColumnData(index, "value", e.target.value)
                    }
                    className="w-full m-2"
                  />
                  <Button
                    color={"light"}
                    onClick={() => deleteColumnData(index)}
                  >
                    Delete
                  </Button>
                </div>
              ))}
              <Button onClick={addColumn} className="m-2">
                Add Header
              </Button>
            </Tabs.Item>
            <Tabs.Item active title="Request Body">
              <Textarea
                className="h-[calc(80vh)] mr-4 resize-none"
                value={requestBody}
                onChange={onChangeRequestBody}
              ></Textarea>
            </Tabs.Item>
          </Tabs>
        </div>

        <div id="responseBlock" className="w-full h-80vh mr-4">
          <Tabs>
            <Tabs.Item active title="Response Body">
              <Textarea
                className="h-[calc(80vh)] mr-4 resize-none"
                value={responseData}
                readOnly
              ></Textarea>
            </Tabs.Item>
            <Tabs.Item active title="Response Header">
              {responseHeader.map((header, index) => (
                <div key={index} className="flex flex-grow p-2">
                  <TextInput
                    type="text"
                    value={header.key}
                    className="w-full m-2"
                  />
                  <TextInput
                    type="text"
                    value={header.value}
                    className="w-full m-2"
                  />
                </div>
              ))}
            </Tabs.Item>
            <Tabs.Item active title="Response Status">
              {responseStatus}
            </Tabs.Item>
          </Tabs>
        </div>
      </div>
    </div>
  );
};
