import './App.css';
import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { toast } from 'react-toastify';
import MaterialReactTable from "material-react-table";
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import 'react-tabs/style/react-tabs.css';

const socket = new WebSocket("ws://127.0.0.1:8010/ws");


function App() {

  const [darkMode, setDarkMode] = useState(localStorage.getItem('darkMode'));

  const [searchKeyword, setSearchKeyword] = useState('');
  const [message, setMessage] = useState('');

  const [containersData, setContainersData] = useState([]);

  const server_url = 'http://localhost:8010';

  useEffect(() => {
    // validatePageAccess();
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

  const handleToggleDarkMode = e => {
    localStorage.setItem('darkMode', darkMode === '1' ? '0' : '1');
    setDarkMode(darkMode === '1' ? '0' : '1');
    toast.info('Not implemented yet')
  }


  const columns = [
    {
      accessorKey: 'name',
      header: 'Name',
    },
    {
      accessorKey: 'value',
      header: 'Value',
    },
  ];


  return (

    <div className="App">


      <div className={"container mt-2"}>
          <img src={"docker_smiley_logo.png"} width={"150"} height={"150"}/>
      </div>


        <div className={"container border p-3"}>


          <Tabs>
            <TabList>
              <Tab>Environments</Tab>
              <Tab>Processes</Tab>
              <Tab>Logs</Tab>
            </TabList>

            <TabPanel>

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
                  columns={columns}
                  data={containersData}
                  enableExpanding
                  getSubRows={(originalRow) => originalRow.subRows} //default, can customize
              />
            </TabPanel>
            <TabPanel>
              <h2>Ps content will be here</h2>
            </TabPanel>
            <TabPanel>
              <h2>Logs content will be here</h2>
            </TabPanel>
          </Tabs>




        </div>

        <hr />

        <div className={'container '}>

          <div className="form-check form-switch d-table">
            <input className="form-check-input" type="checkbox" role="switch" id="flexSwitchCheckDefault"
              onChange={e => handleToggleDarkMode(e)}
              checked={darkMode === '1'}
            />
            <label className="form-check-label float-left" htmlFor="flexSwitchCheckDefault">dark mode</label>
          </div>
          <pre>{message}</pre>

        </div>

    </div>
  );
}

export default App;
