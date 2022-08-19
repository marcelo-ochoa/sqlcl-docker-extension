import { useEffect, useState } from 'react';
import { Box, Grid, LinearProgress, Typography, useMediaQuery } from '@mui/material';
import { createDockerDesktopClient } from '@docker/extension-api-client';


const client = createDockerDesktopClient();

function useDockerDesktopClient() {
  return client;
}

export function App() {
  const [started, setStarted] = useState(false);
  const [ready, setReady] = useState(false);
  const ddClient = useDockerDesktopClient();
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');

  useEffect(() => {
    const start = async (darkMode: boolean) => {
      setStarted(() => false);
      setReady(() => false);

      if (darkMode) {
        await ddClient.extension.vm?.service?.post('/dark', null);
      } else {
        await ddClient.extension.vm?.service?.post('/light', null);
      }

      setStarted(() => true);
    };

    start(prefersDarkMode);
  }, [prefersDarkMode]);

  useEffect(() => {
    if (!started || ready) {
      return;
    }


    const checkIfsqlclIsReady = async () => {
      const result = await ddClient.extension.vm?.service?.get('/ready');
      const ready = Boolean(result);
      if (ready) {
        clearInterval(timer);
      }
      setReady(() => ready);
    };

    let timer = setInterval(() => {
      checkIfsqlclIsReady();
    }, 1000);

    return () => {
      clearInterval(timer);
    };
  }, [started, ready]);

  return (
    <>
      {!ready && (
        <Grid container flex={1} direction="column" padding="16px 32px" height="100%" justifyContent="center" alignItems="center">
          <Grid item>
            <LinearProgress/>
            <Typography mt={2}>
              Waiting for sqlcl to be ready. It may take some seconds if
              it's the first time.
            </Typography>
          </Grid>
        </Grid>
      )}
      {ready && (
        <Box display="flex" flex={1} width="100%" height="100%">
          <iframe src='http://localhost:9890/' width="100%" height="100%" />
        </Box>
      )}
    </>
  );
}
