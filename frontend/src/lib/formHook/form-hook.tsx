import { useRef, useState } from "react";

export const NO_ERR = "__no_err";
export const ERR_EMPTY = "__empty";
export const ERR_TOO_SHORT = "__too_short";
export const ERR_TOO_LONG = "__too_long";
export const ERR_INVALID_LENGTH = "__invalid_length";
export const ERR_INVALID_EMAIL = "__invalid_email";
export const ERR_NO_CAPITALS = "__no_capitals";
export const ERR_NO_DIGITS = "__no_digits";
export const ERR_NO_SPECIAL_CHARACTERS = "__no_special_characters";
export const ERR_MUST_CHECKED = "__must_checked";
export const ERR_INVALID_VALUE = "__invalid_value";
export const ERR_LATIN_ONLY = "__latin_only";
export const ERR_UNKNOWN_FIELD = "__unknown_field";

const SPECIAL_CHARACTERS = /[*\-+!@#$%^&(),.?":{}|<>~\\\/\[\];'=_]/g;

interface FieldText {
  type: "text";
  notEmpty?: boolean;
  latinOnly?: boolean;
  minLength?: number;
  maxLength?: number;
  exactLength?: number;
}

interface FieldEmail {
  type: "email";
}

interface FieldPassword {
  type: "password";
  minLength?: number;
  maxLength?: number;
  includeDigits?: boolean;
  includeCapitals?: boolean;
  includeSpecialCharacters?: boolean;
  latinOnly?: boolean;
}

interface FieldCheckbox {
  type: "checkbox";
  mustChecked?: boolean;
}

type Field = FieldText | FieldEmail | FieldPassword | FieldCheckbox;

const validateFieldText = (value: string, requirements: FieldText) => {
  if ((requirements.notEmpty ?? false) && value.length == 0) return ERR_EMPTY;
  if (value.length < (requirements.minLength ?? 0)) return ERR_TOO_SHORT;
  if (value.length > (requirements.maxLength ?? Infinity)) return ERR_TOO_LONG;

  if (value.length != (requirements.exactLength ?? value.length)) {
    return ERR_INVALID_LENGTH;
  }

  if (requirements.latinOnly && /\W/.test(value)) {
    return ERR_LATIN_ONLY;
  }

  return NO_ERR;
};

const validateFieldEmail = (value: string, requirements: FieldEmail) => {
  // if (!EmailValidator.validate(value)) return ERR_INVALID_EMAIL;
  return NO_ERR;
};

const validateFieldPassword = (value: string, requirements: FieldPassword) => {
  if (value.length < (requirements.minLength ?? 0)) return ERR_TOO_SHORT;
  if (value.length > (requirements.maxLength ?? Infinity)) return ERR_TOO_LONG;

  if ((requirements.includeCapitals ?? false) && value == value.toLowerCase()) {
    return ERR_NO_CAPITALS;
  }

  if ((requirements.includeDigits ?? false) && !/[0-9]/.test(value)) {
    return ERR_NO_DIGITS;
  }

  if (
    (requirements.includeSpecialCharacters ?? false) &&
    !SPECIAL_CHARACTERS.test(value)
  ) {
    return ERR_NO_SPECIAL_CHARACTERS;
  }

  if (
    requirements.latinOnly &&
    /\W/.test(value.replaceAll(SPECIAL_CHARACTERS, ""))
  ) {
    return ERR_LATIN_ONLY;
  }

  return NO_ERR;
};

const validateFieldCheckbox = (value: boolean, requirements: FieldCheckbox) => {
  if ((requirements.mustChecked ?? false) && !value) return ERR_MUST_CHECKED;
  return NO_ERR;
};

const validators = {
  text: validateFieldText,
  email: validateFieldEmail,
  password: validateFieldPassword,
  checkbox: validateFieldCheckbox,
};

type Validator = (value: string | boolean, requirements: Field) => string;

export function useForm(requirements: {
  [key: string]: Field;
}): [
  (values: { [key: string]: string | boolean }) => boolean,
  (field: string) => string,
  (...fields: string[]) => void,
  (next: (...value: any[]) => any) => (...value: any[]) => any
] {
  let initialErrors: { [key: string]: string } = {};

  for (let [field, _] of Object.entries(requirements)) {
    initialErrors[field] = NO_ERR;
  }

  const [errorFields, setErrorFields] = useState(initialErrors);
  const errorRef = useRef(errorFields);
  errorRef.current = errorFields;

  const validate = (values: { [key: string]: string | boolean }) => {
    let noErrors = true;
    const errors: { [key: string]: string } = {};

    for (let [field, value] of Object.entries(values)) {
      if (!Object.hasOwn(requirements, field)) continue;

      const fieldRequirements = requirements[field];
      const validator = validators[fieldRequirements.type] as Validator;

      const err = validator(value, fieldRequirements);

      errors[field] = err;
      if (err != NO_ERR) {
        noErrors = false;
      }
    }

    setErrorFields((prev) => ({ ...prev, ...errors }));

    return noErrors;
  };

  const getError = (field: string) => {
    if (!errorRef.current) return NO_ERR;

    const errors = errorRef.current;

    if (!Object.hasOwn(errors, field)) return ERR_UNKNOWN_FIELD;

    return errors[field];
  };

  const markInvalid = (...fields: string[]) => {
    setErrorFields((prev) => {
      let newErrors = { ...prev };

      for (let field of fields) {
        if (!Object.hasOwn(newErrors, field)) continue;
        newErrors[field] = ERR_INVALID_VALUE;
      }

      return newErrors;
    });
  };

  const withReset = (next: (...value: any[]) => any) => {
    return (...value: any[]) => {
      setErrorFields({ ...initialErrors });
      return next(...value);
    };
  };

  return [validate, getError, markInvalid, withReset];
}
