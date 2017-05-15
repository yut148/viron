import constants from '../../core/constants';

export default {
  show: () => {
    return Promise
      .resolve();
  },

  update: (context, key) => {
    return Promise
      .resolve()
      .then(() => {
        context.commit(constants.MUTATION_CURRENT, key);
      });
  },

  remove: (context) => {
    return Promise
      .resolve()
      .then(() => {
        context.commit(constants.MUTATION_CURRENT, null);
      });
  }
};
