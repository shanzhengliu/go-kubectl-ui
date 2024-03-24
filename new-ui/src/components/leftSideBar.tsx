import {  Sidebar } from 'flowbite-react';

export const  LeftSideBar = (props: {
    show: string;
    setShow: (show: string) => void;
}) => {

return (
    <Sidebar className='w-12' aria-label="Sidebar">
        <Sidebar.Items>
            <Sidebar.ItemGroup>
                <Sidebar.Item  className={props.show === "k8sViewer" ? "text-yellow-700" : ""} onClick={()=>{props.setShow("k8sViewer")}}    >
                  K8S
                </Sidebar.Item>
                <Sidebar.Item className={props.show === "tools" ? "text-yellow-700" : ""} onClick={()=>{props.setShow("tools")}}>
                  Tools
                </Sidebar.Item>
            </Sidebar.ItemGroup>
        </Sidebar.Items>
    </Sidebar>
);
}