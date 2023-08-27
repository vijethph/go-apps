import { useState, useEffect } from "react";
import { Card, Form, Button, ListGroup, Badge, Alert } from "react-bootstrap";

export function Home() {
  const [input, setInput] = useState<number>(0);
  const [seenIndexes, setSeenIndexes] = useState<number[]>([]);
  const [seenValues, setSeenValues] = useState<{
    [key: string]: string;
  }>({});
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    fetchSeenIndexes();
    fetchSeenValues();
  }, []);

  async function fetchSeenIndexes() {
    const response = await fetch("/api/values/all");
    const data = await response.json();
    setSeenIndexes(data);
  }

  async function fetchSeenValues() {
    const response = await fetch("/api/values/current");
    const data = await response.json();
    setSeenValues(data);
  }

  async function calculateFibonacci() {
    setLoading(true);
    await fetch("/api/values", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        index: input,
      }),
    });
    setLoading(false);
    setInput(0);
    fetchSeenValues();
  }

  return (
    <>
      <Card className="mb-3" bg="secondary" text="white">
        <Card.Header>Input Numbers</Card.Header>
        <Card.Body>
          {seenIndexes.map((index) => (
            <Badge bg="primary" className="me-2" key={index}>
              {index}
            </Badge>
          ))}
        </Card.Body>
      </Card>
      <Card>
        <Card.Header>Calculated Fibonacci Numbers</Card.Header>
        <Card.Body>
          <ListGroup horizontal className="d-flex flex-wrap">
            {Object.entries(seenValues).map(([key, value]) => (
              <ListGroup.Item
                className="d-flex justify-content-between mb-2 me-2"
                key={key}
              >
                <Badge bg="primary" className="me-2">
                  {key}
                </Badge>
                <Badge bg="success" className="me-2">
                  {value}
                </Badge>
              </ListGroup.Item>
            ))}
          </ListGroup>
        </Card.Body>
      </Card>
      <Card className="my-3">
        <Card.Header>Calculate Fibonacci</Card.Header>
        <Card.Body>
          <Form onSubmit={(e) => e.preventDefault()}>
            <Form.Group className="mb-3" controlId="formBasicEmail">
              <Form.Label>Enter a number</Form.Label>
              <Form.Control
                type="number"
                placeholder="Enter a number"
                value={input}
                onChange={(e) => setInput(parseInt(e.target.value))}
              />
            </Form.Group>
            <Button variant="primary" onClick={calculateFibonacci}>
              Calculate
            </Button>
          </Form>
          {loading && (
            <Alert variant="warning" className="mt-3">
              Calculating Fibonacci for {input}. Please refresh the page to see
              the result.
            </Alert>
          )}
        </Card.Body>
      </Card>
    </>
  );
}
