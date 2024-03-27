import { leftSideBarStore } from "../react-context/shareContext";
import { Navigator } from "./k8sViewer/navigators";
import { LeftSideBar } from "./leftSideBar";
import { ToolsPanel } from "./tools/toolsPanel";

export const MainView = () => {

const letSideMap: { [key: string]: JSX.Element } = {
  "k8sViewer": <Navigator />,
  "tools": <ToolsPanel/>
};

const show = leftSideBarStore((state) => state.show);
const setShow = leftSideBarStore((state) => state.setShow);

return (
  <div className="h-screen w-screen flex ">
    <div className="w-12 h-full" style={{ position: 'fixed' }}>
      <LeftSideBar show={show} setShow={setShow} />
    </div>
    <div className="flex-grow w-[calc(100vw-3rem)] min-w-0 ml-[3rem]">
      {letSideMap[show] || <div></div>}
    </div>
  </div>
);
};
