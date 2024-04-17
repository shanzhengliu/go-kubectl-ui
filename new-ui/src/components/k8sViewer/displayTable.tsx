import {Button, Checkbox, Table, TextInput} from "flowbite-react";
import React, {useEffect} from "react";
import _ from 'lodash';
import {
    JSXElementConstructor,
    Key,
    ReactElement,
    ReactNode,
    ReactPortal,
} from "react";
import {ContextSwitcher} from "./contextSwitcher";
import {inputHook} from "../../hooks/inputhook";

export function DisplayTable(props: {
    header: any[];
    data: any[][];
    checkbox?: boolean;
    isLoading?:boolean;
    refresh: () => void;
}) {
    const [searchValue, setSearchValue] = React.useState("");
    const [renderData, setRenderData] = React.useState<any[][]>(props.data || []);
    const [originData, setOriginData] = React.useState<any[][]>(props.data || []);
    const [filterTags, setFilterTags] = React.useState<string[]>(
        localStorage.getItem("filterTags")?.split(",") || []
    );
    const [addFilterInput, , onchangeFilterInput] = inputHook("");
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

    useEffect(() => {
        if (filterTags.length !== 0) {
            localStorage.setItem("filterTags", filterTags.join(","));
        }
    }, [filterTags]);

    const addFilterTag = () => {
        if (addFilterInput === "") return;
        setFilterTags([...filterTags, addFilterInput]);

    };

    const clearFlags = () => {
        setFilterTags([]);
        localStorage.removeItem("filterTags");
    }

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
            if (!_.isEqual(renderData, result)) {
                setRenderData(result);
            }
        }
    }, [searchValue, renderData]);

    useEffect(() => {
        setRenderData(props.data);
        setOriginData(props.data);
    }, [props.data]);

    return (
        <div>
            <div className="flex justify-between items-center mb-4">
                <ContextSwitcher onSwitch={props.refresh}/>
                
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
            <div className="flex justify-start items-center m-4">
                Filter Tags:
                <TextInput
                    className="ml-4"
                    value={addFilterInput}
                    onChange={onchangeFilterInput}
                    placeholder="add tag for filter"
                />
                <Button className="ml-4" onClick={addFilterTag}>
                    Add
                </Button>

            </div>


            {filterTags.length === 0 ? null : (
                <div className="flex justify-start items-center m-4">
                    Filters:
                    {filterTags.map((tag) => (
                        <div className="flex">
                            <Button
                                key={tag}
                                color={"success"}
                                onClick={() => {
                                    setSearchValue(tag);
                                }}
                                className="ml-4"
                            >
                                {tag}
                            </Button>

                        </div>
                    ))}
                    <Button onClick={clearFlags} color="failure" className="ml-6">Clear</Button>
                </div>)}

                {props.isLoading &&<div className="flex"><div className="animate-spin rounded-full h-4 w-4 border-b-2 border-gray-900"></div><div>loading</div></div>}
            <div className="w-[calc(100vw-4rem)]">
                <Table hoverable>
                    <Table.Head>
                        {props.checkbox ? (
                            <Table.HeadCell className="p-4">
                                <Checkbox/>
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
                                    <Table.HeadCell className="p-3">
                                        <Checkbox/>
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
        </div>
    );
}
