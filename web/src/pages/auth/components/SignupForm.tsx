import type { HTMLChakraProps } from "@chakra-ui/react";
import {
  Button,
  chakra,
  FormControl,
  FormLabel,
  Input,
  Stack,
  useToast,
} from "@chakra-ui/react";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

import AuthApi from "../../../api/auth-api";
import type { ShortUser } from "../../../models/user";
import { useErrorHandler } from "../../../utils/handle-get-error";

import PasswordField from "./PasswordField";

const SignupForm = (props: HTMLChakraProps<"form">) => {
  const [user, setUser] = useState<ShortUser>({
    email: "",
    password: "",
  });
  const toast = useToast();
  const errorHandler = useErrorHandler();
  const navigate = useNavigate();

  return (
    <chakra.form
      onSubmit={(e) => {
        e.preventDefault();
        e.preventDefault();
        AuthApi.signup(user)
          .then(() => {
            toast({
              title: "Activation link",
              status: "success",
              duration: 1500,
            });
            navigate("/signin");
          })
          .catch(errorHandler);
      }}
      {...props}
    >
      <Stack spacing="6">
        <FormControl id="email">
          <FormLabel>Email address</FormLabel>
          <Input
            value={user.email}
            onChange={(event) =>
              setUser({ ...user, email: event.target.value })
            }
            name="email"
            type="email"
            autoComplete="email"
            required
          />
        </FormControl>
        <PasswordField
          value={user.password}
          onChange={(event) =>
            setUser({ ...user, password: event.target.value })
          }
        />
        <Button type="submit" colorScheme="blue" size="lg" fontSize="md">
          Sign up
        </Button>
      </Stack>
    </chakra.form>
  );
};

export default SignupForm;
