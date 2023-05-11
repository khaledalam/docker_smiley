import './App.css';
import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { toast } from 'react-toastify';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import 'react-tabs/style/react-tabs.css';
import TabPanelEnvironments from "./Components/TabPanelEnvironments";
import TabPanelLogs from "./Components/TabPanelLogs";
import TabPanelProcesses from "./Components/TabPanelProcesses";


function App() {

  const [darkMode, setDarkMode] = useState(localStorage.getItem('darkMode'));

  const [message, setMessage] = useState('');



  useEffect(() => {


  }, []);



  const handleToggleDarkMode = e => {
    localStorage.setItem('darkMode', darkMode === '1' ? '0' : '1');
    setDarkMode(darkMode === '1' ? '0' : '1');
    toast.info('Not implemented yet')
  }



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

              <TabPanelEnvironments />

            </TabPanel>
            <TabPanel>
              <TabPanelProcesses />
            </TabPanel>
            <TabPanel>
              <TabPanelLogs />
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
