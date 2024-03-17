import { Checkbox, Table } from "flowbite-react";
import React from "react";
import {
  JSXElementConstructor,
  Key,
  ReactElement,
  ReactNode,
  ReactPortal,
  useState,
} from "react";

export function DisplayTable(props: {
  header: any[];
  data: any[][];
  checkbox?: boolean;
}) {
  return (
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
          {props.data.map((row, index) => (
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
                    {" "}
                    {React.isValidElement(cell) ? cell : <span>{cell}</span>}
                  </Table.Cell>
                )
              )}
            </Table.Row>
          ))}
        </Table.Body>
      </Table>
    </div>
  );
}
