import type { HTMLChakraProps } from "@chakra-ui/react";
import {
  Button,
  Center,
  chakra,
  FormControl,
  FormLabel,
  Input,
  Stack,
} from "@chakra-ui/react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { JsonErrorResponse } from "utils/json-api-document";

import AuthApi from "../../../api/auth-api";
import { useErrorHandler } from "../../../utils/handle-get-error";
import PasswordField from "./PasswordField";

type ChangePasswordFormProps = HTMLChakraProps<"form"> & { token: string }

const ChangePasswordForm = (props: ChangePasswordFormProps) => {
  const [password, setPassword] = useState("");
  const errorHandler = useErrorHandler();
  const navigate = useNavigate();

  return (
    <chakra.form
      onSubmit={(e) => {
        e.preventDefault();
        AuthApi.changePassword(props.token, password)
          .catch((e: JsonErrorResponse) => {
            if (
              e.statusText.toLowerCase().includes('token') 
              || e.data.error.title.toLowerCase().includes('token')
            ) {
              const invalidTokenErr = {
                status: 401,
                statusText: 'Invalid token',
                data: {
                  error: {
                    title: 'Please, try to reset email again',
                  }
                }
              } as JsonErrorResponse
              errorHandler(invalidTokenErr)
            } else {
              errorHandler(e)
            }
          })
          .finally(() => navigate('/signin', { replace: true }));
      }}
      {...props}
    >
      <Stack spacing="6">
        <PasswordField
          value={password}
          onChange={(event) => setPassword(event.target.value)}
        />
        <Button type="submit" colorScheme="blue" size="lg" fontSize="md">
          Change password
        </Button>
      </Stack>
    </chakra.form>
  );
};
export default ChangePasswordForm;
