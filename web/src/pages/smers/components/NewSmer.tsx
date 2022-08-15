import { AddIcon, DeleteIcon } from "@chakra-ui/icons";
import {
  Box,
  Button,
  Center,
  HStack,
  IconButton,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Textarea,
  useToast,
} from "@chakra-ui/react";
import { useEffect, useMemo, useReducer, useRef, useState } from "react";

import SmersApi from "../../../api/smers-api";
import Step from "../../../components/stepper/Step";
import type { SmerItemEdit } from "../../../models/smer";
import { NewSmerDto, StepState } from "../../../models/smer";
import { useErrorHandler } from "../../../utils/handle-get-error";

interface NewSmerProps {
  isOpen: boolean;
  isLoading: boolean;
  activeSmerId: undefined | number;
  onClose: () => void;
  onSmerSave: () => void;
}

enum StepAction {
  Forward,
  Backward,
  Target,
}

const NewSmer = ({
  onClose,
  onSmerSave,
  isOpen,
  activeSmerId,
}: NewSmerProps) => {
  const [smer, setSmer] = useState(
    new NewSmerDto({
      situation: "",
      thoughts: [],
      emotions: [],
      reactions: [],
    }).newState
  );

  const [newItem, setNewItem] = useState("");
  const newItemRef = useRef<HTMLInputElement>(null);

  const errorHandler = useErrorHandler();
  const [isLoading, setIsLoading] = useState(false);
  const toast = useToast();

  useEffect(() => {
    if (activeSmerId) {
      setIsLoading(true);
      SmersApi.view(activeSmerId)
        .then((payload) => setSmer(new NewSmerDto(payload).newState))
        .catch(errorHandler)
        .finally(() => setIsLoading(false));
    }
  }, [activeSmerId]);

  const nextStepReducer = (
    state: StepState,
    action: { type: StepAction; payload?: StepState }
  ) => {
    if (action.type === StepAction.Forward) {
      switch (state) {
        case StepState.Situation:
          return StepState.Thoughts;
        case StepState.Thoughts:
          return StepState.Emotions;
        case StepState.Emotions:
          return StepState.Reactions;
        default:
          return StepState.Situation;
      }
    } else if (action.type === StepAction.Backward) {
      switch (state) {
        case StepState.Thoughts:
          return StepState.Situation;
        case StepState.Emotions:
          return StepState.Thoughts;
        case StepState.Reactions:
          return StepState.Emotions;
        default:
          return StepState.Situation;
      }
    } else {
      return action.payload || StepState.Situation;
    }
  };

  const [currentStep, nextStep] = useReducer(
    nextStepReducer,
    StepState.Situation
  );

  const clearState = () => {
    setNewItem("");
    setSmer(
      new NewSmerDto({
        situation: "",
        thoughts: [],
        emotions: [],
        reactions: [],
      }).newState
    );
    nextStep({ type: StepAction.Target, payload: StepState.Situation });
  };

  const onActionDone = () => {
    toast({
      title: "Success",
      status: "success",
      duration: 1500,
    });
    clearState();
    onSmerSave();
  };

  const onSave = () => {
    setIsLoading(true);
    if (activeSmerId) {
      SmersApi.update(activeSmerId, NewSmerDto.toSmerDto(smer))
        .then(() => {
          onActionDone();
        })
        .catch(errorHandler)
        .finally(() => setIsLoading(false));
    } else {
      SmersApi.create(NewSmerDto.toSmerDto(smer))
        .then(() => {
          onActionDone();
        })
        .catch(errorHandler)
        .finally(() => setIsLoading(false));
    }
  };

  const nextBtnIsDisabled = useMemo(() => {
    switch (currentStep) {
      case StepState.Situation:
        return !smer.situation.length;
      case StepState.Thoughts:
        return !smer.thoughts.length;
      case StepState.Emotions:
        return !smer.emotions.length;
      case StepState.Reactions:
        return !smer.reactions.length;
      default:
        return false;
    }
  }, [currentStep, smer]);

  const saveBtnIsDisabled = useMemo(() => {
    return (
      !smer.situation.length ||
      !smer.thoughts.length ||
      !smer.emotions.length ||
      !smer.reactions.length
    );
  }, [smer]);

  useEffect(() => newItemRef.current?.focus(), [currentStep]);

  const addItem = () => {
    setSmer(NewSmerDto.addItem(smer, currentStep, newItem));
    setNewItem("");
    newItemRef.current?.focus();
  };

  const removeItem = (id: string) => {
    setSmer(NewSmerDto.removeItem(smer, currentStep, id));
  };

  const setItem = (value: string, id: string) => {
    setSmer(NewSmerDto.setItem(smer, currentStep, value, id));
  };

  const setSituation = (value: string) => {
    setSmer(NewSmerDto.setSituation(smer, value));
  };

  const onDialogClose = () => {
    clearState();
    onClose();
  };

  return (
    <Modal
      onClose={onDialogClose}
      size="lg"
      isOpen={isOpen}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent maxH="calc(100% - 120px)">
        <ModalHeader>New SMER</ModalHeader>
        <ModalCloseButton />
        <ModalBody overflowY="scroll">
          <Box mx="auto" maxW="3xl" py="10" px={{ base: "6", md: "8" }}>
            <nav aria-label="Progress steps">
              <HStack as="ol" listStyleType="none" spacing="0">
                <Step
                  isCurrent={currentStep === StepState.Situation}
                  isDone={!!smer.situation.length}
                  onClick={() =>
                    nextStep({
                      type: StepAction.Target,
                      payload: StepState.Situation,
                    })
                  }
                >
                  Situation
                </Step>
                <Step
                  isCurrent={currentStep === StepState.Thoughts}
                  isDone={!!smer.thoughts.length}
                  disabled={!smer.situation.length}
                  onClick={() =>
                    nextStep({
                      type: StepAction.Target,
                      payload: StepState.Thoughts,
                    })
                  }
                >
                  Thoughts
                </Step>
                <Step
                  isCurrent={currentStep === StepState.Emotions}
                  isDone={!!smer.emotions.length}
                  disabled={!smer.thoughts.length}
                  onClick={() =>
                    nextStep({
                      type: StepAction.Target,
                      payload: StepState.Emotions,
                    })
                  }
                >
                  Emotions
                </Step>
                <Step
                  isCurrent={currentStep === StepState.Reactions}
                  isDone={!!smer.reactions.length}
                  disabled={!smer.emotions.length}
                  onClick={() =>
                    nextStep({
                      type: StepAction.Target,
                      payload: StepState.Reactions,
                    })
                  }
                >
                  Reaction
                </Step>
              </HStack>
            </nav>
            {currentStep !== StepState.Situation && (
              <Box>
                {smer[currentStep].map((stepField: SmerItemEdit) => (
                  <Center key={stepField.id} my={2} display="flex">
                    <Input
                      value={stepField.value}
                      onChange={(event) =>
                        setItem(event.target.value, stepField.id)
                      }
                      color="teal"
                      placeholder="Item"
                      _placeholder={{ color: "inherit" }}
                    />
                    <IconButton
                      my={2}
                      ml={2}
                      icon={<DeleteIcon />}
                      aria-label="Remove item"
                      onClick={() => removeItem(stepField.id)}
                    />
                  </Center>
                ))}
                <Center display="flex">
                  <Input
                    value={newItem}
                    ref={newItemRef}
                    onChange={(event) => setNewItem(event.target.value)}
                    color="info"
                    placeholder="New item"
                    _placeholder={{ color: "inherit" }}
                  />
                  <IconButton
                    my={2}
                    ml={2}
                    icon={<AddIcon />}
                    disabled={!newItem.length}
                    aria-label="Add item"
                    onClick={addItem}
                  />
                </Center>
              </Box>
            )}
            {currentStep === StepState.Situation && (
              <Box>
                <Textarea
                  isInvalid={!smer.situation.length}
                  value={smer.situation}
                  my={2}
                  onChange={(event) => setSituation(event.target.value)}
                  color="teal"
                  placeholder="Situation"
                  _placeholder={{ color: "inherit" }}
                />
              </Box>
            )}
          </Box>
        </ModalBody>
        <ModalFooter justifyContent="space-between">
          <Button
            isLoading={isLoading}
            disabled={currentStep === StepState.Situation}
            onClick={() => nextStep({ type: StepAction.Backward })}
          >
            Back
          </Button>
          {currentStep !== StepState.Reactions && (
            <Button
              isLoading={isLoading}
              disabled={nextBtnIsDisabled}
              onClick={() => nextStep({ type: StepAction.Forward })}
            >
              Next
            </Button>
          )}
          {currentStep === StepState.Reactions && (
            <Button
              isLoading={isLoading}
              disabled={saveBtnIsDisabled}
              onClick={onSave}
            >
              Save
            </Button>
          )}
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default NewSmer;
