import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';

class App extends Component {
  constructor() {
    super();
    this.state = {
      echoResult: "",
    };
  }


  componentDidMount() {
  
    fetch('/api/echo?data=testeroo')
    .then(results => results.text())
    .then(data => {
      console.log("attempted fetch");
      console.log(data);
      this.setState({ echoResult: data });
    })
    .catch(error => console.log(error));
  }

  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <p>
            Edit <code>src/App.js</code> and save to reload.
            echo results: {this.state.echoResult}
          </p>
          <a
            className="App-link"
            href="https://reactjs.org"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn React
          </a>
        </header>
      </div>
    );
  }
}

export default App;
