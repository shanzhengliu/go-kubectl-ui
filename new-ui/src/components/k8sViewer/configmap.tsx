import { useEffect, useState } from "react";
import { DisplayTable } from "./displayTable";
import { authVerify, axiosInstance } from "../../utils/axios";
import { CONFIGMAP, CONFIGMAP_DETAIL } from "../../utils/endpoints";
import { Button, Modal } from "flowbite-react";
import hljs from "highlight.js";
export function Configmap() {
  const [tableData, setTableData] = useState<any[][]>([]);
  const [openModal, setOpenModal] = useState(false);
  const [modelHeader, setModelHeader] = useState("");
  const [modalData, setModalData] = useState<{ [key: string]: any }>({});
  const dataFecth = async () => {
    const response = await authVerify();
    if (response == "error") {
      return;
    }
    axiosInstance
      .get(CONFIGMAP, {
        data: {},
        headers: {
          "Content-Type": "application/json",
        },
      })
      .then((response) => {
        const responseData: any[][] = [];
        response.data.map((item: any) => {
          responseData.push([
            item.name,
            item.namespace,
            <Button
              onClick={() => {
                setOpenModal(true);
                setModelHeader(item.name);
                axiosInstance
                  .get(`${CONFIGMAP_DETAIL}?configmap=${item.name}`)
                  .then((res) => {
                    setModalData(res.data);
                  });
              }}
            >
              View
            </Button>,
          ]);
        });
        setTableData(responseData);
      });
  };
  hljs.highlightAll();
  useEffect(() => {
    dataFecth();
  }, []);

  return (
    <div>
      <DisplayTable
        data={tableData}
        header={["Configmap", "Namespace", ""]}
        refresh={dataFecth}
      />
      <div>
        <div>
          <Modal show={openModal} onClose={() => setOpenModal(false)}>
            <Modal.Header>{modelHeader}</Modal.Header>
            <Modal.Body>
              <table className="table-auto w-full">
                <thead>
                  <tr>
                    <th className="px-4 py-2">Key</th>
                    <th className="px-4 py-2">Value</th>
                  </tr>
                </thead>
                <tbody>
                  {Object.entries(modalData).map(([key, value], index) => {
                    return (
                      <tr key={index} className="border-t">
                        <td className="px-4 py-2 font-medium">{key}</td>
                        <td className="px-4 py-2">
                          <pre className="bg-gray-100 p-2 rounded">
                            <code className="language-json">{value}</code>
                          </pre>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </Modal.Body>
            <Modal.Footer>
              <Button onClick={() => setOpenModal(false)}>OK</Button>
            </Modal.Footer>
          </Modal>
        </div>
      </div>
    </div>
  );
}
