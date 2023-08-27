import { Button, Card } from "react-bootstrap";
import { Link } from "react-router-dom";

export function OtherPage() {
  return (
    <Card className="mb-3">
      <Card.Header>Other Page</Card.Header>
      <Card.Body>
        <p>
          Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed
          malesuada, sapien vel bibendum bibendum, sapien sapien bibendum
          sapien, vel bibendum sapien sapien sit amet sapien.
        </p>
      </Card.Body>
      <Card.Footer>
        <Link to="/">
          <Button>Go back to home page</Button>
        </Link>
      </Card.Footer>
    </Card>
  );
}
