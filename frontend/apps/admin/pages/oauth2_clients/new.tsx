import { useRouter } from 'next/router';
import { useForm, zodResolver } from '@mantine/form';
import { TextInput, Button, Group, Container } from '@mantine/core';
import { z } from 'zod';

import { OAuth2Client, OAuth2ClientCreationRequestInput } from '@dinnerdonebetter/models';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { AppLayout } from '../../src/layouts';

const oauth2ClientCreationFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
});

export default function OAuth2ClientCreator(): JSX.Element {
  const router = useRouter();

  const creationForm = useForm({
    initialValues: {
      name: '',
      description: '',
    },
    validate: zodResolver(oauth2ClientCreationFormSchema),
  });

  const submit = async () => {
    const validation = creationForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const submission = new OAuth2ClientCreationRequestInput({
      name: creationForm.values.name,
      description: creationForm.values.description,
    });

    const apiClient = buildLocalClient();

    await apiClient
      .createOAuth2Client(submission)
      .then((result: OAuth2Client) => {
        if (result) {
          router.push(`/oauth2_clients/${result.id}`);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="Create New OAuth2 Client">
      <Container size="sm">
        <form onSubmit={creationForm.onSubmit(submit)}>
          <TextInput label="Name" placeholder="thing" {...creationForm.getInputProps('name')} />
          <TextInput label="Description" placeholder="thing" {...creationForm.getInputProps('description')} />

          <Group position="center">
            <Button type="submit" mt="sm" fullWidth>
              Submit
            </Button>
          </Group>
        </form>
      </Container>
    </AppLayout>
  );
}
