import { StrictMode } from 'react';
import { Button, Container, Title } from '@mantine/core';

import { AppLayout } from '../src/layouts';
import axios from 'axios';

export default function Web(): JSX.Element {
  return (
    <StrictMode>
      <AppLayout title="">
        <Container size="xs">
          <Title order={5}>Welcome</Title>

          <Button
            onClick={() => {
              axios
                .get('/api/v1/valid_ingredients')
                .then((response) => {
                  console.log(response);
                })
                .catch((error) => {
                  console.error(error);
                });
            }}
          >
            Click Me
          </Button>
        </Container>
      </AppLayout>
    </StrictMode>
  );
}
