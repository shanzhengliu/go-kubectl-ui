import { Button, Label, Modal, Select, TextInput } from "flowbite-react";
import { useEffect, useState } from "react";
import { axiosInstance } from "../utils/axios";
import { CONTEXT_CHANGE, CONTEXT_LIST, CURRENT_CONTEXT } from "../utils/endpoints";

export function ContextSwitcher(props: { onSwitch: () => void }) {
  const [currentContext, setCurrentContext] = useState("NA");
  const [contextSelected, setContextSelected] = useState("");
  const [currentNamespace, setCurrentNamespace] = useState("NA");
  const [inputNamespace, setInputNamespace] = useState("");
  const [contextList, setContextList] = useState<any[]>([]);
  const [openModal, setOpenModal] = useState(false);

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
        props.onSwitch();
        setOpenModal(false);
      });
  };

  return (
    <>
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
    </>
  );
}