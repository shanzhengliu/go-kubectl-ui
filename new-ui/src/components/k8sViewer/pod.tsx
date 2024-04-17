import {useEffect, useState} from "react";
import {DisplayTable} from "./displayTable";
import {authVerify, axiosInstance} from "../../utils/axios";
import Swal from "sweetalert2";
import {POD, PODLOGS, PODYAML, ROLLING_LOG, WEBSHELL} from "../../utils/endpoints";
import {Button} from "flowbite-react";

export function Pod() {
    const [renderData, setRenderData] = useState<any[][]>([]);
    const [loading, setLoading] = useState(false);
    const dataFetch = async () => {
        const response = await authVerify();
        if (response == "error") {
            return;
        }
        setLoading(true);
        axiosInstance
            .get(POD, {
                data: {},
                headers: {
                    "Content-Type": "application/json",
                },
            })
            .then((response) => {
                const responseData: any[][] = [];
                for (let i = 0; i < response.data.length; i++) {
                    for (let j = 0; j < response.data[i].images.length; j++) {
                        responseData.push([
                            response.data[i].name,
                            response.data[i].namespace,
                            response.data[i].images[j].containerName,
                            response.data[i].images[j].name,
                            <ImageId id={response.data[i].images[j].id}></ImageId>,
                            <StatusRunning
                                status={response.data[i].images[j].containerStatus}
                            />,
                            response.data[i].createTime,
                            <div className={"flex pl-2"}>
                                <div className={"mr-1"}>
                                    <RollingLog
                                        pod={response.data[i].name}
                                        container={response.data[i].images[j].containerName}
                                    />
                                </div>
                                <div className={"mr-1"}>
                                    <Logs

                                        pod={response.data[i].name}
                                        container={response.data[i].images[j].containerName}
                                    /></div>
                                <div className={"mr-1"}>
                                    <Yaml pod={response.data[i].name}/>
                                </div>
                                <WebShell
                                    pod={response.data[i].name}
                                    container={response.data[i].images[j].containerName}
                                /></div>
                        ]);
                    }
                }
                setRenderData(responseData);
                setLoading(false);
            });
    };
    useEffect(() => {
        dataFetch();
    }, []);
    return (
        <div>
            <DisplayTable
                isLoading={loading}
                data={renderData}
                refresh={dataFetch}
                header={[
                    "Pod",
                    "Namespace",
                    "Container",
                    "Image",
                    "Image Id",
                    "Status",
                    "Create Time",
                    "",
                ]}
            />
        </div>
    );
}

function Logs(props: { pod: string; container: string }) {
    return (
        <Button gradientMonochrome="info">
            <div>
                <a
                    href={`${PODLOGS}?pod=${props.pod}&container=${props.container}`}
                    target="_blank"
                >
                    Logs
                </a>
            </div>
        </Button>
    );
}

function StatusRunning(props: { status: string }) {
    switch (props.status) {
        case "Running":
            return <span className="text-green-500">Running</span>;
        case "Pending":
            return <span className="text-red-500">Pending</span>;
        case "Waiting":
            return <span className="text-yellow-500">Waiting</span>;
        default:
            return <span className="text-blue-500">{props.status}</span>;
    }
}

function WebShell(props: { pod: string; container: string }) {
    return (
        <div>
            <Button gradientMonochrome="success">
                <a
                    href={`${WEBSHELL}?pod=${props.pod}&container=${props.container}`}
                    target="_blank"
                >
                    Shell
                </a>
            </Button>
        </div>
    );
}

function RollingLog(props: { pod: string; container: string }) {
    return (
        <div>
            <Button className={"text-white"} gradientDuoTone="pinkToOrange">
                <a
                    href={`${ROLLING_LOG}?pod=${props.pod}&container=${props.container}`}
                    target="_blank"
                    className={"break-after-all whitespace-nowrap"}
                >
                    Rolling Logs
                </a>
            </Button>
        </div>
    );
}

function ImageId(props: { id: string }) {
    return (
        <div>
            <Button
                color="light"
                onClick={() => {
                    Swal.fire({
                        title: "Image ID",
                        text: props.id,
                        confirmButtonText: "ok",
                    });
                }}
            >
                ID
            </Button>
        </div>
    );
}

function Yaml(props: { pod: string }) {
    return (
        <div>
            <Button gradientDuoTone="cyanToBlue">
                <a href={`${PODYAML}?pod=${props.pod}`} target="_blank">
                    Yaml
                </a>
            </Button>
        </div>
    );
}
