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

import AuthApi from "../../../api/auth-api";
import { useErrorHandler } from "../../../utils/handle-get-error";
import PasswordField from "./PasswordField";

const ChangePasswordForm = (props: HTMLChakraProps<"form"> & { token: string }) => {
  const [password, setPassword] = useState<string>("");
  const errorHandler = useErrorHandler();
  const navigate = useNavigate();

  return (
    <chakra.form
      onSubmit={(e) => {
        e.preventDefault();
        AuthApi.changePassword(props.token, password)
          .then(() => {
            navigate('/signin', { replace: true });
          })  
          .catch(errorHandler);
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
