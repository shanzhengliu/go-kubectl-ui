import { Label, List } from "flowbite-react";
import { FaFile, FaFolder, FaFolderOpen } from "react-icons/fa";
import { OPENAPI_HELPER_GET_FILES } from "../../utils/endpoints";
import { useEffect, useState } from "react";


interface FileTree {
    [key: string]: boolean | FileTree;
  }
  
export const FileTreeComponent = (props:{
    setSelectFile?: (path: string) => void
}) => {
    const [fileTree, setFileTree] = useState<FileTree | null>(null);
    const [expandedDirs, setExpandedDirs] = useState<string[]>([]);
  
    useEffect(() => {
      fetch(OPENAPI_HELPER_GET_FILES)
        .then(response => response.json())
        .then(data => setFileTree(data))
        .catch(error => console.error('Error fetching file tree:', error));
    }, []);
  
    const handleClick = (path: string) => {
        setExpandedDirs(prevExpandedDirs => {
          if (prevExpandedDirs.includes(path)) {
            return prevExpandedDirs.filter(dir => dir !== path);
          } else {
            return [...prevExpandedDirs, path];
          }
        });
      };
    
      const handleFileClick = (path: string) => {
        console.log(`Clicked file: ${path}`);
        if (props.setSelectFile) props.setSelectFile(path);
      };
  
    const renderFileTree = (tree: FileTree, parentPath: string = '',level: number = 0): JSX.Element[] => {
      return Object.entries(tree).map(([key, value]) => {
        const path = `${parentPath}/${key}`;
        const isDirectory = typeof value === 'object';
        const isExpanded = expandedDirs.includes(path);
  
        return (
          <List.Item
            key={path}
           
            className={`level-${level}`}
            style={{ paddingLeft: `${level * 20}px`,listStyle: 'none', cursor: 'pointer' }}
           
          >
            {isDirectory ? (
            <>
              <div onClick={() => handleClick(path)} className="flex" >
                {isExpanded ? <FaFolderOpen className="mr-2" /> : <FaFolder className="mr-2" />}
                {key}
              </div>
              {isExpanded && (
                <List className="ml-4" style={{ listStyle: 'none' }}>
                  {renderFileTree(value as FileTree, path, level + 1)}
                </List>
              )}
            </>
          ) : (
            <div onClick={() => handleFileClick(path)} className="flex">
              <FaFile className="mr-2" />
              {key}
            </div>
          )}
          </List.Item>
        );
      });
    };
  
    return (
      <div className="mt-2"   >
         <h2>File Tree</h2>
        {fileTree && (
          <List style={{ listStyle: 'none' }}>
            {renderFileTree(fileTree)}
          </List>
        )}
      </div>
    );
  };