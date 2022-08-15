import type { HTMLChakraProps } from "@chakra-ui/react";
import {
  Button,
  Center,
  chakra,
  Container,
  FormControl,
  FormLabel,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Stack,
  Avatar as UserAvatar,
  useDisclosure,
} from "@chakra-ui/react";
import { useState, useEffect } from "react";
import Avatar from "react-avatar-edit";

import UsersApi from "../../api/users-api";
import type { User } from "../../models/user";
import { useErrorHandler } from "../../utils/handle-get-error";
import { useSuccessHandler } from "../../utils/handle-success";

const Profile = (props: HTMLChakraProps<"form">) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [user, setUser] = useState<User>({
    avatar: null,
    email: "",
    name: "",
    surname: "",
    username: "",
  });
  const [preview, setPreview] = useState<string | undefined>("");
  const errorHandler = useErrorHandler();
  const successHandler = useSuccessHandler("Success");

  const onCloseImage = () => {
    setPreview(undefined);
  };

  const onCropImage = (newPreview: string) => {
    setPreview(newPreview);
  };

  const onConfirmImage = () => {
    setUser({ ...user, avatar: preview, avatarId: null });
    onClose();
  };

  const fetchUser = () => {
    UsersApi.view().then(setUser).catch(errorHandler);
  };

  const updateUser = () => {
    UsersApi.update(user)
      .then(() => {
        successHandler("User updated successfully");
        fetchUser();
      })
      .catch(errorHandler);
  };

  useEffect(fetchUser, []);

  return (
    <Container>
      <chakra.form
        onSubmit={(e) => {
          e.preventDefault();
          updateUser();
        }}
        {...props}
      >
        <Stack spacing="6">
          <UserAvatar
            width={200}
            height={200}
            src={
              user.avatarId
                ? `${process.env.SERVER_URL}/${user.avatar}`
                : user.avatar
            }
            onClick={onOpen}
          />
          <Modal isOpen={isOpen} onClose={onClose}>
            <ModalOverlay />
            <ModalContent>
              <ModalHeader>Pick image</ModalHeader>
              <ModalCloseButton />
              <ModalBody>
                <Center>
                  <Avatar
                    width={400}
                    height={250}
                    onCrop={onCropImage}
                    onClose={onCloseImage}
                    src={
                      user.avatarId && user.avatar
                        ? `${process.env.SERVER_URL}/${user.avatar}`
                        : user.avatar
                    }
                  />
                </Center>
              </ModalBody>
              <ModalFooter justifyContent="space-between">
                <Button variant="ghost" onClick={onClose}>
                  Close
                </Button>
                <Button colorScheme="blue" onClick={onConfirmImage}>
                  Confirm
                </Button>
              </ModalFooter>
            </ModalContent>
          </Modal>
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
          <FormControl id="name">
            <FormLabel>Name</FormLabel>
            <Input
              value={user.name}
              onChange={(event) =>
                setUser({ ...user, name: event.target.value })
              }
              name="name"
              type="name"
              autoComplete="name"
            />
          </FormControl>
          <FormControl id="name">
            <FormLabel>Surname</FormLabel>
            <Input
              value={user.surname}
              onChange={(event) =>
                setUser({ ...user, surname: event.target.value })
              }
              name="surname"
              type="surname"
              autoComplete="surname"
            />
          </FormControl>
          <FormControl id="username">
            <FormLabel>Username</FormLabel>
            <Input
              value={user.username}
              onChange={(event) =>
                setUser({ ...user, username: event.target.value })
              }
              name="username"
              type="username"
              autoComplete="username"
            />
          </FormControl>
          <Center />
          <Button type="submit" colorScheme="blue" size="lg" fontSize="md">
            Update
          </Button>
        </Stack>
      </chakra.form>
    </Container>
  );
};
export default Profile;
