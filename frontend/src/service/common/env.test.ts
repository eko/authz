import { isEnabled } from './env';

test.each([
  [undefined, false],
  ['', false],
  ['false', false],
  ['0', false],
  ['no', false],
  ['true', true],
  ['TRUE', true],
  ['1', true],
])('isEnabled(%s) returns %s', (value, expected) => {
  expect(isEnabled(value)).toBe(expected);
});
