# Example from https://github.com/mlr-org/mlr3gallery
# MIT License Copyright (c) 2019 mlr-org

library(mlr3learners)

# creates mlr3 task from scratch, from a data.frame
# 'target' names the column in the dataset we want to learn to predict
task = as_task_classif(iris, target = "Species")
# in this case we could also take the iris example from mlr3's dictionary of shipped example tasks
# 2 equivalent calls to create a task. The second is just sugar for the user.
task = mlr_tasks$get("iris")
task = tsk("iris")
# print(task)
# create learner from dictionary of mlr3learners
# 2 equivalent calls:
learner_1 = mlr_learners$get("classif.rpart")
learner_1 = lrn("classif.rpart")
# print(learner_1)

# train learner on subset of task
learner_1$train(task, row_ids = 1:120)
# this is what the decision tree looks like
# print(learner_1$model)
# predict using observations from task
prediction = learner_1$predict(task, row_ids = 121:150)
# predict using "new" observations from an external data.frame
prediction = learner_1$predict_newdata(newdata = iris[121:150, ])
# print(prediction)

# head(as.data.table(mlr_measures))
scores = prediction$score(msr("classif.acc"))
print(scores)
scores = prediction$score(msrs(c("classif.acc", "classif.ce")))
print(scores)
cm = prediction$confusion
print(cm)
