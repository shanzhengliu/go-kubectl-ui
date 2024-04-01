import {useEffect, useState} from "react";
import {DisplayTable} from "./displayTable";
import {authVerify, axiosInstance} from "../../utils/axios";
import {SERVICE, START_PORT_FORWARD, STOP_PORT_FORWARD} from "../../utils/endpoints";
import {Button, TextInput} from "flowbite-react";
import {errorHander} from "../../utils/tools.tsx";

export function Service() {
    const [renderData, setRenderData] = useState<any[][]>([]);
    const [localPorts, setLocalPorts] = useState<any>([]);
    const [servicePorts, setServicePorts] = useState<any>([]);
    const [response, setResponse] = useState<any>([]);


    useEffect(() => {
        const responseData: any[] = [];
        if (response.length == 0 || localPorts.length == 0 || servicePorts.length == 0) {
            return;
        }
        for (let i = 0; i < response.length; i++) {
            responseData.push([
                response[i].name,
                response[i].namespace,
                response[i].type,
                response[i].selector,
                <TextInput readOnly={!!response[i].isForward} value={localPorts[i]}
                           onChange={(e) => handleLocalPortChange(i, e.target.value)}/>,
                <TextInput readOnly={!!response[i].isForward} value={servicePorts[i]}
                           onChange={(e) => handleServicePortChange(i, e.target.value)}/>,
                response[i].isForward ? <Button color={"failure"}
                                                onClick={() => stopForward(i, response[i].namespace, response[i].name)}>Stop</Button> :
                    <Button color={"success"}
                            onClick={() => startForward(i, response[i].namespace, response[i].name)}>Start</Button>,
            ]);
        }
        setRenderData(responseData);

    }, [response, localPorts, servicePorts]);


    const stopForward = (index: number, namespace: string, service: string) => {
        axiosInstance.get(STOP_PORT_FORWARD, {

            headers: {
                "Content-Type": "application/json",
            },
            params: {
                localPort: localPorts[index],
                servicePort: servicePorts[index],
                namespace: namespace,
                service: service,
            },
        }).then(async () => {
            await fetchData();
        });
    }

    const startForward = (index: number, namespace: string, service: string) => {
        axiosInstance.get(START_PORT_FORWARD, {

            headers: {
                "Content-Type": "application/json",
            },
            params: {
                localPort: localPorts[index],
                servicePort: servicePorts[index],
                namespace: namespace,
                service: service,
            },
        }).then(async () => {
            await fetchData();
        }).catch((error) => {
            errorHander(error)
        });
    }

    const handleServicePortChange = (index: number, value: string) => {
        const newServicePorts = [...servicePorts];
        newServicePorts[index] = value;
        setServicePorts(newServicePorts);
    };
    const handleLocalPortChange = (index: number, value: string) => {
        const newLocalPorts = [...localPorts];
        newLocalPorts[index] = value;
        setLocalPorts(newLocalPorts);
    }

    const fetchData = async () => {
        const response = await authVerify();
        if (response == "error") {
            return;
        }
        axiosInstance
            .get(SERVICE, {
                data: {},
                headers: {
                    "Content-Type": "application/json",
                },
            })
            .then((response): void => {
                const localPortsArray = response.data.map((item: { localPort: any; }) => item.localPort || "");
                const servicePortsArray = response.data.map((item: { servicePort: any; }) => item.servicePort || "");
                setLocalPorts(localPortsArray);
                setServicePorts(servicePortsArray);
                setResponse(response.data);
            });
    }

    useEffect( () => {
        fetchData();
    }, []);


    return (
        <div>
            <DisplayTable
                header={["Service", "NameSpace", "Type", "Selector", "Local Port", "Service Port", ""]}
                data={renderData}
                refresh={fetchData}
            ></DisplayTable>
        </div>
    );
}
