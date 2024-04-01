import {Button, Label, Select, Textarea, TextInput} from "flowbite-react";
import {inputHook} from "../../hooks/inputhook.tsx";
import {useEffect, useState} from "react";
import {axiosInstance} from "../../utils/axios.tsx";


export const OpenAPIMockView = (props: { mockLoader: any, mockPort: string }) => {
    const [path, setPath, onChangePath] = inputHook("");
    const [method, setMethod, onChangeMethod] = inputHook("");
    const [status, setStatus, onChangeStatus] = inputHook("");
    const [contentType, setContentType, onChangeContentType] = inputHook("");
    const [methodList, setMethodList] = useState<string[]>([]);
    const [statusList, setStatusList] = useState<string[]>([]);
    const [contentTypeList, setContentTypeList] = useState<string[]>([]);
    const [examples, setExamples, onChangeExamples] = inputHook("");
    const [exampleList, setExampleList] = useState<string[]>([]);


    const [finalUrl, setFinalUrl, onChangefinalUrl] = inputHook("");
    const [finalMethod, setFinalMethod, onChangefinalMethod] = inputHook("");
    const [finalStatus, setFinalStatus, onChangefinalStatus] = inputHook("");
    const [finalContentType, setFinalContentType, onChangefinalContentType] = inputHook("");
    const [finalExamples, setFinalExamples, onChangefinalExamples] = inputHook("");
    const [finalBody, setFinalBody, onChangefinalBody] = inputHook("{}");

    const [responseBody, setResponseBody] = useState<any>({});
    const [responseStatus, setResponseStatus] = useState<number>(0);


    const resetStatus = () => {
        setStatus("");
        setStatusList([]);
    }

    const resetMethod = () => {
        setMethod("");
        setMethodList([]);

    }

    const resetContentType = () => {
        setContentType("");
        setContentTypeList([]);
    }
    const resetPath = () => {
        setPath("");
    }

    const resetExamples = () => {
        setExamples("");
        setExampleList([]);

    }


    useEffect(() => {
        if (props.mockLoader) {
            resetStatus();
            resetMethod();
            resetPath();
            resetContentType()
            resetExamples()
            setFinalUrl("");
            setFinalExamples("");
            setFinalContentType("");
            setFinalMethod("");
            setFinalStatus("");
            setFinalBody("{}")
            setResponseBody({})
            setResponseStatus(0)

        }
    }, [props.mockLoader]);

    useEffect(() => {
        if (path) {

            resetStatus();
            resetMethod();
            resetContentType()
            resetExamples()
            setMethodList(Object.keys(props.mockLoader.paths[path]));
        }

    }, [path]);


    useEffect(() => {
        if (path && method) {
            resetStatus();
            resetContentType()
            resetExamples()
            setStatusList(Object.keys(props.mockLoader.paths[path][method].responses));
        }
    }, [method]);

    useEffect(() => {
        if (path && method && status) {
            resetContentType()
            resetExamples()
            setContentTypeList(Object.keys(props.mockLoader.paths[path][method].responses[status].content));
        }
    }, [status]);

    useEffect(() => {
        if (path && method && status && contentType) {
            if (props.mockLoader.paths[path][method].responses[status].content[contentType].examples) {
                setExampleList(Object.keys(props.mockLoader.paths[path][method].responses[status].content[contentType].examples));
            }

        }
    }, [contentType]);

    const generateRequestForm = () => {
        if (path && method && status && contentType) {
            const finalUrl = "http://localhost:"+props.mockPort + path
            const finalMethod = method;
            const finalStatus = status;
            const finalContentType = contentType;
            const finalExamples = examples;
            setFinalUrl(finalUrl);
            setFinalMethod(finalMethod);
            setFinalStatus(finalStatus);
            setFinalContentType(finalContentType);
            setFinalExamples(finalExamples);
        }
    }

    const sendRequest = () => {
        console.log(finalUrl, finalMethod, finalStatus, finalContentType, finalExamples);
        const options = {
            method: finalMethod,
            url: finalUrl,
            headers: {
                "openapi-status-code": finalStatus,
                "openapi-content-type": finalContentType,
                "openapi-examples": finalExamples,
                "Content-Type": finalContentType,
                "Accept": "*/*",
            },
            data: finalBody
        }
        axiosInstance.request(options).then(response => {
            setResponseBody(response.data);
            setResponseStatus(response.status);
        }).catch(error => {
            setResponseBody(error.response.data);
            setResponseStatus(error.response.status);

        })
    }


    return (<div className={"w-full flex"}>

                <div className={"flex w-full p-2"}>

                    <div className="p-2">
                        <h2>Request Details Select</h2>
                        <Label value="Path"/>
                        <Select value={path} onChange={onChangePath}>
                            <option value={""}>Select Path</option>
                            {props.mockLoader.paths && Object.keys(props.mockLoader.paths).map((item, index) => {
                                    return <option value={item} key={index}>{item}</option>
                                }
                            )}
                        </Select>
                        <Label value={"Method"}/>
                        <Select value={method} onChange={onChangeMethod}>
                            <option value={""}>Select Method</option>
                            {methodList && methodList.map((item, index) => {
                                return <option value={item} key={index}>{item}</option>
                            })}
                        </Select>
                        <Label value={"Status"}/>
                        <Select value={status} onChange={onChangeStatus}>
                            <option value={""}>Select Status</option>
                            {statusList && statusList.map((item, index) => {
                                return <option value={item} key={index}>{item}</option>
                            })}
                        </Select>
                        <Label value={"Content Type"}/>
                        <Select value={contentType} onChange={onChangeContentType}>
                            <option value={""}>Select Content Type</option>
                            {contentTypeList && contentTypeList.map((item, index) => {
                                    return <option value={item} key={index}>{item}</option>
                                }
                            )}
                        </Select>
                        <Label value={"Examples"}/>
                        <Select value={examples} onChange={onChangeExamples}>
                            <option value={""}>Select Examples</option>
                            {exampleList && exampleList.map((item, index) => {
                                return <option value={item} key={index}>{item}</option>
                            })}
                        </Select>
                        <Button className={"mt-2"} onClick={generateRequestForm}>Generate Form</Button>
                    </div>
                </div>
            <div className={"w-full p-2"}>
                <h2>Generate Request Detail</h2>
                <div>
                    <Label value={"URL"}/>
                    <TextInput value={finalUrl} onChange={onChangefinalUrl}></TextInput>
                </div>
                <div>
                    <Label value={"Method"}/>
                    <TextInput value={finalMethod} onChange={onChangefinalMethod}></TextInput>
                </div>
                <div>
                    <Label value={"Status"}/>
                    <TextInput value={finalStatus} onChange={onChangefinalStatus}></TextInput>
                </div>
                <div>
                    <Label value={"Content Type"}/>
                    <TextInput value={finalContentType} onChange={onChangefinalContentType}></TextInput>
                </div>
                <div>
                    <Label value={"Example"}/>
                    <TextInput value={finalExamples} onChange={onChangefinalExamples}></TextInput>
                </div>
                <div>
                    <Label value={"Body"}/>
                    <Textarea value={finalBody} onChange={onChangefinalBody}></Textarea>
                </div>
                <Button  className={"mt-2"} onClick={sendRequest}>Send</Button>
            </div>
            <div className={"w-full h-screen p-2"}>
                <h2>Response</h2>
                <Label value={"Status"}/>
                <TextInput  readOnly className={""} value={responseStatus}></TextInput>
                <Label className={"mt-2"} value={"Body"}/>
                <Textarea  readOnly className={"resize-none h-screen"} value= {responseBody && JSON.stringify(responseBody,null,2 )}></Textarea>

            </div>
        </div>
    )
}