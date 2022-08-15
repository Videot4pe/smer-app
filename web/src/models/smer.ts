import { v4 as uuidv4 } from "uuid";

export interface SmerDto {
  id?: number;
  userId?: number;

  situation: string;
  thoughts: string[];
  emotions: string[];
  reactions: string[];

  createdAt?: string;
  updatedAt?: string;
}

export enum StepState {
  Situation = "situation",
  Thoughts = "thoughts",
  Emotions = "emotions",
  Reactions = "reactions",
}

export interface SmerItemEdit {
  value: string;
  id: string;
}

export class NewSmerDto {
  public state: Record<string, any> = {
    situation: "",
    thoughts: [],
    emotions: [],
    reactions: [],
  };

  constructor(smer: SmerDto) {
    this.state.situation = smer.situation;
    this.state.thoughts = smer.thoughts.map((thought) => ({
      value: thought,
      id: uuidv4(),
    }));
    this.state.emotions = smer.emotions.map((emotion) => ({
      value: emotion,
      id: uuidv4(),
    }));
    this.state.reactions = smer.reactions.map((reaction) => ({
      value: reaction,
      id: uuidv4(),
    }));
  }

  get newState() {
    return { ...this.state };
  }

  static addItem(smer: Record<string, any>, key: string, value: string) {
    const newSmer = { ...smer };
    newSmer[key].push({ value, id: uuidv4() });
    return newSmer;
  }

  static setItem(
    smer: Record<string, any>,
    key: string,
    value: string,
    id: string
  ) {
    const newSmer = { ...smer };
    const element = newSmer[key].find((el: SmerItemEdit) => el.id === id);
    const index = newSmer[key].indexOf(element);
    newSmer[key][index].value = value;
    return newSmer;
  }

  static removeItem(smer: Record<string, any>, key: string, id: string) {
    const newSmer = { ...smer };
    const element = newSmer[key].find((el: SmerItemEdit) => el.id === id);
    const index = newSmer[key].indexOf(element);
    newSmer[key].splice(index, 1);
    return newSmer;
  }

  static setSituation(smer: Record<string, any>, value: string) {
    const newSmer = { ...smer };
    newSmer.situation = value;
    return newSmer;
  }

  static toSmerDto(smer: Record<string, any>): SmerDto {
    return {
      situation: smer.situation,
      thoughts: smer.thoughts.map((thought: SmerItemEdit) => thought.value),
      emotions: smer.emotions.map((emotion: SmerItemEdit) => emotion.value),
      reactions: smer.reactions.map((reaction: SmerItemEdit) => reaction.value),
    };
  }
}
