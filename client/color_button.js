const root = document.documentElement;

function App() {
  const [bckgrnd, setBckgrnd] = React.useState('#46924d');

  root.style.setProperty('--bckgrnd-color', bckgrnd);

  async function changeColor() {
    const hexClr = await getHexColor();
    setBckgrnd(hexClr);
  }

  return (
    <div>
      <p>The current color is {bckgrnd}.</p>

      {/* We can have different types of events here that are supported by React: https://reactarmory.com/guides/react-events-cheatsheet */}
      <Button text="Change Color" onClick={changeColor} />
      <AddForm />
    </div>
  );
}

async function getHexColor() {
  try {
    const response = await axios.get('http://localhost:8080');
    return `#${response.data}`;
  } catch (error) {
    console.error(error);
  }
}

const Button = props => {
  return (
    <button onClick={props.onClick} className="color-btn">
      {props.text}
    </button>
  );
};

const AddForm = () => {
  const [addends, setAddends] = React.useState({
    num1: 0,
    num2: 0,
  });
  const [result, setResult] = React.useState(0);

  function handleChange(event) {
    setAddends({
      ...addends,
      [event.target.name]: event.target.value,
    });
  }

  function handleSubmit(event) {
    add();
    event.preventDefault();
  }

  async function add() {
    let bodyFormData = new FormData();
    bodyFormData.set('num1', addends.num1);
    bodyFormData.set('num2', addends.num2);
    try {
      const response = await axios({
        method: 'post',
        url: 'http://localhost:8080',
        data: bodyFormData,
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      console.log(response.data);
      setResult(response.data);
    } catch (error) {
      console.error(error);
    }
  }

  return (
    <React.Fragment>
      <form onSubmit={handleSubmit} method="POST" target="_parent">
        <label htmlFor="num1">Number One</label>
        <br />
        <input type="number" id="num1" name="num1" onChange={handleChange} />
        <br />
        <label htmlFor="num2">Number Two</label>
        <br />
        <input type="number" id="num2" name="num2" onChange={handleChange} />
        <br />
        <input type="submit" value="Add" />
      </form>
      <Results addends={addends} result={result} />
    </React.Fragment>
  );
};

const Results = props => {
  const { addends, result } = props;
  const { num1, num2 } = addends;
  return (
    <p>
      The result of {num1} + {num2} is {result}.
    </p>
  );
};

function renderApp() {
  ReactDOM.render(<App />, document.getElementById('react-app'));
}

renderApp();
