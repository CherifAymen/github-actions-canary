import React, {Component} from 'react';
import './App.css';
import Polarity from "./components/Polarity";
import { Container, Row, Col, Button,InputGroup, InputGroupAddon, Input, Jumbotron} from 'reactstrap';



class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            val:'',
            sentence: '',
            polarity: undefined
        };
        this.handleChange = this.handleChange.bind(this);
    };

    analyzeSentence() {

        if (this.state.val !== ''){
            this.state.sentence = this.state.val 
            this.state.val = '' 
    
            fetch('http://localhost:8080/sentiment', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({sentence: this.state.sentence})
            })
                .then(response => response.json())
                .then(data => this.setState(data));
        }
    }



    handleChange({ target }) {
        this.setState({
          [target.name]: target.value
        });
      }

      onEnterPress = e => {
        if (e.key === 'Enter') {
            this.analyzeSentence();
        }
    };

    render() {
        const polarityComponent = this.state.polarity !== undefined ?
            <Polarity sentence={this.state.sentence} polarity={this.state.polarity}/> :
            null;

        return (
            <div className="deco">
                <Container>
                    <Row>
                        <Col xs="6" className="my-5 mx-auto deco2 py-5">
                        <div>
                            <Jumbotron >
                                <h1 className="display-3 text-center mb-3">Analyseur de sentiment</h1>
                               
                                <InputGroup>
                                    <Input placeholder="Envoyer une phrase" name="val" value={ this.state.val} onChange={ this.handleChange} onKeyUp={this.onEnterPress.bind(this)}/>
                                    <InputGroupAddon addonType="append"><Button color="secondary" onClick={this.analyzeSentence.bind(this)}>Envoyer</Button></InputGroupAddon>
                                </InputGroup>
                                        
                                <hr className="my-2" />
                                <ul>
                                    <li> Le résultat doit être entre -1 et 1.</li>
                                    <li> Le résultat s'approchera de 1 si le sens de la phrase est appréciatif </li>
                                    <li>Le résultat s'approchera de -1 si le sens de la phrase est péjoratif</li>
                                    <li>Le résultat sera 0 si la phrase est neutre</li>
                                    <li>Les phrases doivent être en anglais</li>
                                </ul>
                               <div className="my-4">
                                    {polarityComponent}
                               </div>

                            </Jumbotron>
                        </div>
                            
                        </Col>
                        <Col className="center my-5 py-5" xs="6">
                                <img   width="150%" src="back.png" alt="Card image cap" />
                        </Col>
                    </Row>
                    <div className="mb-5 pb-5 pt-3 mt-3 ">
                        <h3 className="text-center text-secondary "> &#169; Aymen Cherif </h3>
                    </div>
                </Container>
            </div>
        );
    }
}

export default App;
