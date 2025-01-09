import { GetServerSideProps, GetServerSidePropsResult } from 'next';
import { Container } from '@mantine/core';

import React from 'react';
import { AppLayout } from '../src/layouts';

declare interface HomePageProps {
  NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID: string;
  NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET: string;
  NEXT_PUBLIC_SEGMENT_API_TOKEN: string;
  REWRITE_COOKIE_HOST_FROM: string;
  REWRITE_COOKIE_HOST_TO: string;
  REWRITE_COOKIE_SECURE: string;
  NEXT_API_ENDPOINT: string;
  NEXT_PUBLIC_API_ENDPOINT: string;
  NEXT_COOKIE_ENCRYPTION_KEY: string;
  NEXT_BASE64_COOKIE_ENCRYPT_IV: string;
}

export const getServerSideProps: GetServerSideProps = async (): Promise<GetServerSidePropsResult<HomePageProps>> => {
  return {
    props: {
      NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID: process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID || '',
      NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET: process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET || '',
      NEXT_PUBLIC_SEGMENT_API_TOKEN: process.env.NEXT_PUBLIC_SEGMENT_API_TOKEN || '',
      REWRITE_COOKIE_HOST_FROM: process.env.REWRITE_COOKIE_HOST_FROM || '',
      REWRITE_COOKIE_HOST_TO: process.env.REWRITE_COOKIE_HOST_TO || '',
      REWRITE_COOKIE_SECURE: process.env.REWRITE_COOKIE_SECURE || '',
      NEXT_API_ENDPOINT: process.env.NEXT_API_ENDPOINT || '',
      NEXT_PUBLIC_API_ENDPOINT: process.env.NEXT_PUBLIC_API_ENDPOINT || '',
      NEXT_COOKIE_ENCRYPTION_KEY: process.env.NEXT_COOKIE_ENCRYPTION_KEY || '',
      NEXT_BASE64_COOKIE_ENCRYPT_IV: process.env.NEXT_BASE64_COOKIE_ENCRYPT_IV || '',
    },
  };
};

function HomePage(props: HomePageProps) {
  return (
    <AppLayout title="Meal Plans" userLoggedIn>
      <Container size="xs">
        <hr />
        <p>server-side</p>
        <hr />

        <p>NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID: {`"${props.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID}"`}</p>
        <p>NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET: {`"${props.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET}"`}</p>
        <p>NEXT_PUBLIC_SEGMENT_API_TOKEN: {`"${props.NEXT_PUBLIC_SEGMENT_API_TOKEN}"`}</p>
        <p>REWRITE_COOKIE_HOST_FROM: {`"${props.REWRITE_COOKIE_HOST_FROM}"`}</p>
        <p>REWRITE_COOKIE_HOST_TO: {`"${props.REWRITE_COOKIE_HOST_TO}"`}</p>
        <p>REWRITE_COOKIE_SECURE: {`"${props.REWRITE_COOKIE_SECURE}"`}</p>
        <p>NEXT_API_ENDPOINT: {`"${props.NEXT_API_ENDPOINT}"`}</p>
        <p>NEXT_PUBLIC_API_ENDPOINT: {`"${props.NEXT_PUBLIC_API_ENDPOINT}"`}</p>
        <p>NEXT_COOKIE_ENCRYPTION_KEY: {`"${props.NEXT_COOKIE_ENCRYPTION_KEY}"`}</p>
        <p>NEXT_BASE64_COOKIE_ENCRYPT_IV: {`"${props.NEXT_BASE64_COOKIE_ENCRYPT_IV}"`}</p>

        <hr />
        <p>browser-side</p>
        <hr />

        <p>NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID: {`"${process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_ID}"`}</p>
        <p>
          NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET:{' '}
          {`"${process.env.NEXT_DINNER_DONE_BETTER_OAUTH2_CLIENT_SECRET}"`}
        </p>
        <p>NEXT_PUBLIC_SEGMENT_API_TOKEN: {`"${process.env.NEXT_PUBLIC_SEGMENT_API_TOKEN}"`}</p>
        <p>REWRITE_COOKIE_HOST_FROM: {`"${process.env.REWRITE_COOKIE_HOST_FROM}"`}</p>
        <p>REWRITE_COOKIE_HOST_TO: {`"${process.env.REWRITE_COOKIE_HOST_TO}"`}</p>
        <p>REWRITE_COOKIE_SECURE: {`"${process.env.REWRITE_COOKIE_SECURE}"`}</p>
        <p>NEXT_API_ENDPOINT: {`"${process.env.NEXT_API_ENDPOINT}"`}</p>
        <p>NEXT_PUBLIC_API_ENDPOINT: {`"${process.env.NEXT_PUBLIC_API_ENDPOINT}"`}</p>
        <p>NEXT_COOKIE_ENCRYPTION_KEY: {`"${process.env.NEXT_COOKIE_ENCRYPTION_KEY}"`}</p>
        <p>NEXT_BASE64_COOKIE_ENCRYPT_IV: {`"${process.env.NEXT_BASE64_COOKIE_ENCRYPT_IV}"`}</p>
      </Container>
    </AppLayout>
  );
}

export default HomePage;
