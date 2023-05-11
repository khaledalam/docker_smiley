import '../App.css';
import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { toast } from 'react-toastify';
import MaterialReactTable from "material-react-table";
import 'react-tabs/style/react-tabs.css';

const socket = new WebSocket("ws://127.0.0.1:8010/ws");


function TabPanelEnvironments() {


  const [searchKeyword, setSearchKeyword] = useState('');
  const [message, setMessage] = useState('');

  const [containersData, setContainersData] = useState([]);

  const server_url = 'http://localhost:8010';

  useEffect(() => {
    fetchContainers();

    socket.onopen = () => {
      setMessage('WS Connected')
    };

    socket.onmessage = (e) => {
      setMessage("Get message from server: " + e.data)
    };

    return () => {
      socket.close()
    }

  }, []);


  const fetchContainers = async () => {

    await axios.get(server_url + "/env/list", {}).then(res => {

      setContainersData(Object.keys(res?.data)?.map((containerKey, index) => {

        return {
          name: containerKey,
          subRows: res?.data[containerKey].map((data, index) => {
            return {
              name: data?.Names[0],
              value: data?.State,
              subRows: Object.values(data?.Envs).map((env, index) => {
                return {
                  name: env?.Name,
                  value: env?.Value,
                  subRows: env?.Levels.map((level, index) => {
                    return {
                      name: level?.LevelString,
                      value: level.Value
                    }
                  })
                }
              })
            }
          })
        }
      }));

    }).catch(err => {
      toast.error("fail to fetch containers");
    });

  };



  const handleKeyDown = async (e) => {
    if (e.key === 'Enter') {
      // handle search action
    }
  }


  const handleChangeSearchKeywordInput = e => {
    setSearchKeyword(e.target.value);
  }



  const envColumns = [
    {
      accessorKey: 'name',
      header: 'Name',
    },
    {
      accessorKey: 'value',
      header: 'Value',
    },
  ];


  return <>
          <div className="input-group">
            <span className="input-group-text" id="basic-addon1">Env key or value</span>

            <input type="text"
                   disabled={true}
                   required={true}
                   value={searchKeyword}
                   onChange={e => handleChangeSearchKeywordInput(e)}
                   onKeyDown={e => handleKeyDown(e)}
                   className="form-control rounded"
                   placeholder="e.g. PATH" />
            <button type="button" className="btn btn-outline-primary"
                    onClick={() => null}
                    disabled={searchKeyword.length < 2}

            >Search</button>


          </div>


          <MaterialReactTable
              columns={envColumns}
              data={containersData}
              enableExpanding
              getSubRows={(originalRow) => originalRow.subRows} //default, can customize
          />
  </>;
}

export default TabPanelEnvironments;
