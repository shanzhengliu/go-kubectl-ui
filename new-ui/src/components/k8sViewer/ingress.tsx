import {useEffect, useState} from "react";
import {authVerify, axiosInstance} from "../../utils/axios";
import {INGRESS} from "../../utils/endpoints";
import {DisplayTable} from "./displayTable";

export function Ingress() {
    const [renderData, setRenderData] = useState<any[]>([]);
    const [loading, setLoading] = useState(false);
    const fetchData = async () => {
        const response = await authVerify();
        if (response == "error") {
            return;
        }
        setLoading(true);
        axiosInstance
            .get(INGRESS, {
                data: {},
                headers: {
                    "Content-Type": "application/json",
                },
            })
            .then((response) => {
                const responseData: any[] = [];
                for (let i = 0; i < response.data.length; i++) {
                    for (let j = 0; j < response.data[i].rules.length; j++) {
                        responseData.push([
                            response.data[i].name,
                            response.data[i].namespace,
                            <a target="_blank" className="underline text-blue-700"
                               href={"https://" + response.data[i].rules[j].host}>{"https://" + response.data[i].rules[j].host}</a>,
                        ]);
                    }
                }
                setRenderData(responseData);
                setLoading(false);
            });
    };
    useEffect(() => {
        fetchData();
    }, []);

    return <div>
        <DisplayTable header={["Ingress", "Namespace", "Host"]} data={renderData} refresh={fetchData}
                      isLoading={loading}></DisplayTable>
    </div>;
}