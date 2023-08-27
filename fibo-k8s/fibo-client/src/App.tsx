import { BrowserRouter as Router, Routes, Route, Link } from "react-router-dom";
import { Navbar, Nav, Container } from "react-bootstrap";
import { Home } from "./components/Home";
import { OtherPage } from "./components/OtherPage";

import "./App.css";

function App() {
  return (
    <Router>
      <Navbar bg="primary" expand="lg" data-bs-theme="dark">
        <Container>
          <Navbar.Brand href="#home">Fibonacci Calculator</Navbar.Brand>
          <Navbar.Toggle aria-controls="basic-navbar-nav" />
          <Navbar.Collapse id="basic-navbar-nav">
            <Nav className="ms-auto">
              <Link to="/" className="nav-link">
                Home
              </Link>
              <Link to="/others" className="nav-link">
                Other Page
              </Link>
            </Nav>
          </Navbar.Collapse>
        </Container>
      </Navbar>
      <Container className="my-3">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/others" element={<OtherPage />} />
        </Routes>
      </Container>
    </Router>
  );
}

export default App;
