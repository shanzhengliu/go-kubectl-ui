import { useState } from 'react';
import { DockerShell } from './dockerShell';
import { JsonFormatter } from './jsonFormatter';
import { HttpHelper } from './httpHelper';
import { EncodeHelper } from './encodeHelper';
import { OpenAPIOnline } from './openAPIOnline';

export const ToolsPanel = () => {
  const toolsMap: { [key: string]: any } = {
    dockerShell: { title: 'Docker Shell', component: <DockerShell iframeKey={Date.now()} /> },
    jsonFormatter: { title: 'Json Formatter', component: <JsonFormatter /> },
    httpHelper: { title: 'Http Helper', component: <HttpHelper /> },
    encodeHelper: { title: 'Encode Helper', component: <EncodeHelper /> },
    openAPIOnline: { title: 'OpenAPI Online', component: <OpenAPIOnline /> },
  };

  const [currentTool, setCurrentTool] = useState('default');
  const [isExpanded, setIsExpanded] = useState(true);

  return (
    <div className="flex">
    
      <div className={`fixed min-h-screen ${isExpanded ? 'w-64' : 'w-16'} transition-all duration-300 ease-in-out bg-gray-200 text-gray-800 flex flex-col items-center py-2 shadow-lg`} >
        <button onClick={() => setIsExpanded(!isExpanded)} className="m-4 p-2 rounded hover:bg-gray-300 focus:outline-none">
          {isExpanded ? '<<' : '>>'}
        </button>
        {isExpanded && (
          Object.keys(toolsMap).map((key) => (
            <CardComponent
              key={key}
              title={toolsMap[key].title}
              onClick={() => {
                setCurrentTool(key);
                setIsExpanded(false);
              }}
            />
          ))
        )}
      </div>
      <div 
        className="flex-grow p-4 transition-margin duration-300 ease-in-out" 
        style={{ marginLeft: isExpanded ? '256px' : '64px' }}
      >
        {toolsMap[currentTool] ? toolsMap[currentTool].component : 'Select a tool from the sidebar'}
      </div>
    </div>
  );
};

function CardComponent({ title, onClick }: { title: string; onClick: () => void }) {
  return (
    <div className="cursor-pointer py-2 px-4 w-full text-center hover:bg-gray-300 rounded transition-colors duration-200" onClick={onClick}>
      <h5 className="text-sm">{title}</h5>
    </div>
  );
}
