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
import { useAtom } from "jotai";
import { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";

import AuthApi from "../../../api/auth-api";
import type { ShortUser } from "../../../models/user";
import { jwtToken, refreshJwtToken } from "../../../store";
import { useErrorHandler } from "../../../utils/handle-get-error";
import PasswordField from "./PasswordField";

const ChangePasswordForm = (props: HTMLChakraProps<"form"> & { hash: string }) => {
  const [password, setPassword] = useState<string>("");
  const errorHandler = useErrorHandler();
  const navigate = useNavigate();
  const location = useLocation();
  // @ts-ignore
  const from = location.state?.from || "/singin";

  return (
    <chakra.form
      onSubmit={(e) => {
        e.preventDefault();
        AuthApi.changePassword(props.hash, password)
          .then(() => {
            navigate(from, { replace: true });
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
